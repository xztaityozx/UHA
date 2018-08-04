//Copyright © 2018 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

const (
	RESERVE string = "reserve"
	DONE    string = "done"
	FAILED  string = "failed"
	SRUN    string = "srun"
	RUN     string = "run"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "シミュレーションを実行します",
	Long:  `シミュレーションセットを実行します`,
	Run: func(cmd *cobra.Command, args []string) {
		num, _ := cmd.PersistentFlags().GetInt("number")
		all, _ := cmd.PersistentFlags().GetBool("all")

		msg := SlackMessage{
			StartTime:  time.Now(),
			Failed:     -1,
			Succsess:   -1,
			SubMessage: "成功、失敗の個数はシミュレーションセットの個数です",
		}

		tasks, err := readRunTasks(num, all)
		if err != nil {
			PostFailed(config.SlackConfig, err)
			log.Fatal(err)
		}

		summ, err := Run(&tasks)
		if err != nil {
			PostFailed(config.SlackConfig, err)
			log.Fatal(err)
		}

		msg.FinishedTime = time.Now()
		cnt := 0
		for _, v := range summ {
			if v.Status {
				cnt++
			}
		}

		printSummary(&summ)

		msg.Succsess = cnt
		msg.Failed = len(tasks) - cnt

		Post(config.SlackConfig, msg)
	},
}

type RunTask struct {
	Simulation Simulation
	Addfile    string
	SPI        []string
	Dst        []string
	ResultFile []string
	Base       string
	SEED       int
	ACE        string
	TaskFile   string
}

// シミュレーションの書き出し先を作る
func tryMkRunDstDir(srt *RunTask) error {
	base := filepath.Join(srt.Simulation.DstDir, fmt.Sprintf("VtpVolt%.4f_VtnVolt%.4f", srt.Simulation.Vtp.Voltage, srt.Simulation.Vtn.Voltage))
	srt.Base = base
	for i, v := range srt.Simulation.Monte {
		srt.Dst = append(srt.Dst, filepath.Join(base, fmt.Sprintf("SEED%03d/Monte%s", srt.SEED, v)))
		rltdir := filepath.Join(base, fmt.Sprintf("SEED%03d/Result/", srt.SEED))
		srt.ResultFile = append(srt.ResultFile, filepath.Join(rltdir, fmt.Sprintf("%s.csv", v)))

		// Directory作る
		if err := tryMkdir(srt.Dst[i]); err != nil {
			return err
		}

		if err := tryMkdir(rltdir); err != nil {
			return err
		}
	}
	return nil
}

//Run用のAddfileを作る
func tryMkRunAddfile(srt *RunTask) error {
	dir := filepath.Join(srt.Base, "Addfiles")
	if err := tryMkdir(dir); err != nil {
		return err
	}

	f := filepath.Join(dir, fmt.Sprintf("RunAddfileSEED%03d.txt", srt.SEED))
	srt.Addfile = f

	var addfile string
	src := filepath.Join(ConfigDir, "addfile.txt")
	if tmplt, err := ioutil.ReadFile(src); err != nil {
		return err
	} else {
		addfile = fmt.Sprintf(string(tmplt), srt.SEED)
	}

	return ioutil.WriteFile(f, []byte(addfile), 0644)
}

// Run用のSPIを作る
func tryMkRunSPI(srt *RunTask) error {

	for _, v := range srt.Simulation.Monte {
		// SPIの名前、長すぎ
		f := filepath.Join(srt.Simulation.SimDir, fmt.Sprintf("VtpVolt%.4f_VtnVolt%.4f_SEED%03d_Monte%s.spi",
			srt.Simulation.Vtp.Voltage,
			srt.Simulation.Vtn.Voltage,
			srt.SEED,
			v))

		srt.SPI = append(srt.SPI, f)

		// 書き込み
		if spi, err := getSPIScript(srt.Simulation, v, srt.Addfile); err != nil {
			return err
		} else {
			if e := ioutil.WriteFile(f, spi, 0644); e != nil {
				return e
			}
		}
	}
	return nil
}

// XMLをコピーする
func tryCopyRunXmls(srt RunTask) error {
	for i, v := range srt.Dst {
		// resultsMap.xml
		{
			src := filepath.Join(SelfPath, "templates", "resultsMap.xml")
			dst := filepath.Join(v, "resultsMap.xml")
			if res, err := ioutil.ReadFile(src); err != nil {
				return err
			} else {
				if e := ioutil.WriteFile(dst, res, 0644); e != nil {
					return e
				}
			}
		}
		// results.xml
		{
			src := filepath.Join(SelfPath, "templates", srt.Simulation.Monte[i])
			dst := filepath.Join(v, "results.xml")
			if res, err := ioutil.ReadFile(src); err != nil {
				return err
			} else {
				if e := ioutil.WriteFile(dst, res, 0644); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

// ごみ処理
func removeRunGarbage(srt RunTask) error {
	// SPI
	for _, v := range srt.SPI {
		if err := os.Remove(v); err != nil {
			return err
		}
	}

	// Dst
	for _, v := range srt.Dst {
		if err := os.RemoveAll(v); err != nil {
			return err
		}
	}

	return nil
}

// ACEを生成して書き込む
func tryMkRunACE(srt *RunTask) error {
	f := filepath.Join(srt.Base, "extract.ace")
	srt.ACE = f
	ace := getACEScript(srt.Simulation.Signal, srt.Simulation.Range)

	return ioutil.WriteFile(f, ace, 0644)
}

// シミュレーションコマンドを生成する
func mkRunCommand(srt RunTask) ([]string, error) {

	if len(srt.Simulation.Monte) == 0 {
		return []string{}, errors.New("タスクが不正です(Monteが0個です)")
	}

	var rt []string
	for i, _ := range srt.Simulation.Monte {
		dst := srt.Dst[i]
		spi := srt.SPI[i]

		command := fmt.Sprintf("cd %s && hspice -mt 2 -i %s -o ./hspice &> ./hspice.log && wv -k -ace_no_gui %s", dst, spi, srt.ACE)
		rt = append(rt, command)
	}

	return rt, nil
}

// ファイルからタスクを解釈する
func MakeRunTask(f string) (RunTask, error) {
	log.Println("Make Run Task from : ", f)
	var s Simulation
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return RunTask{}, err
	}

	if err := json.Unmarshal(b, &s); err != nil {
		return RunTask{}, err
	}

	log.Println(s)

	var rt RunTask = RunTask{
		Simulation: s,
		SEED:       s.SEED,
		TaskFile:   f,
	}
	return rt, nil
}

// Reserveからnum個または全部を読み出す
func readRunTasks(num int, all bool) ([]RunTask, error) {
	files, err := ioutil.ReadDir(ReserveRunDir)
	if err != nil {
		return nil, err
	}

	// numが存在する個数以上、もしくはallが有効なら全部を指定する
	if num > len(files) || all {
		num = len(files)
	}

	if num == 0 {
		return nil, errors.New("タスクファイルがありません")
	}

	var rt []RunTask

	for i, v := range files {
		if v.IsDir() {
			continue
		}
		res, err := MakeRunTask(filepath.Join(ReserveRunDir, v.Name()))
		if err != nil {
			return nil, err
		}

		rt = append(rt, res)
		if i == num {
			break
		}
	}

	return rt, nil
}

// Runする
func Run(tasks *[]RunTask) ([]SRunSummary, error) {
	var rt []SRunSummary
	for _, task := range *tasks {
		res, err := run(task)
		if err != nil {
			// 継続する
			if ContinueWhenFaild {
				log.Println(err)
				continue
			}
			return nil, err
		}
		rt = append(rt, res)
		if err := sendTaskFile(task.TaskFile, res.Status); err != nil {
			if ContinueWhenFaild {
				log.Println(err)
			} else {
				log.Fatal(err)
			}
		}
	}

	return rt, nil

}

func run(task RunTask) (SRunSummary, error) {
	var summary SRunSummary = SRunSummary{
		Name: fmt.Sprintf("VtpVolt:%.4f_VtnVolt:%.4f Sigma:%.4f SEED:%03d",
			task.Simulation.Vtn.Voltage,
			task.Simulation.Vtp.Voltage,
			task.Simulation.Vtn.Sigma,
			task.SEED),
		StartTime: time.Now(),
		Status:    false,
	}
	// ディレクトリ作る
	if err := tryMkRunDstDir(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	log.Println("MkDir :", task.Base)

	// ACE
	if err := tryMkRunACE(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	log.Println("Write ACE to ", task.ACE)

	// Addfile
	if err := tryMkRunAddfile(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	log.Println("Write Addfile to ", task.Addfile)

	// SPI
	if err := tryMkRunSPI(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	log.Println("Write SPI to ", task.SPI[0])

	// XMLs
	if err := tryCopyRunXmls(task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	log.Println("Write Xmls")

	// コマンド生成
	commands, err := mkRunCommand(task)
	if err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	log.Fatal(commands[0])

	log.Println("Start Simulation ", task.TaskFile)
	// spinner
	spin := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	spin.Suffix = "Running... "
	spin.FinalMSG = "Simulation Set had finished!!\n"
	spin.Start()
	defer spin.Stop()

	// waitgroup
	var wg sync.WaitGroup

	// シミュレーション開始
	for i, command := range commands {
		go func(command string, cnt int, l int) {
			wg.Add(1)
			defer wg.Done()
			err := exec.Command("bash", "-c", command).Run()
			if err != nil {
				summary.FinishTime = time.Now()
				if ContinueWhenFaild {
					log.Println(err)
					return
				}
				log.Fatal(err)
			}

			log.Printf("Simulation finished(%d/%d)\n", i, l)

			if err := ExtractFromStoreCSV(task.Dst[i], task.ResultFile[i]); err != nil {
				summary.FinishTime = time.Now()
				if ContinueWhenFaild {
					log.Println(err)
					return
				}
				log.Fatal(err)
			}

			log.Printf("Extract Complete(%d/%d)\n", i, l)

		}(command, i, len(commands))
	}

	wg.Wait()

	summary.Status = true
	summary.FinishTime = time.Now()

	if GargabeCollect {
		if err := removeRunGarbage(task); err != nil {
			log.Println("後片付けに失敗しました")
		}
	}

	return summary, nil

}

// Store.csvからファイルに書き出す
func ExtractFromStoreCSV(SrcDir string, DstFile string) error {
	src := filepath.Join(SrcDir, "store.csv")
	dst := DstFile

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	origin := strings.Split(string(b), "\n")
	var store []string
	// 整形1
	for _, line := range origin {
		// #で始まる行とTIME,SIGNAL行を飛ばす
		if len(line) == 0 || line[0] == '#' || line[0] == 'T' {
			continue
		}

		// 2カラム目を取り出す
		t := strings.Split(strings.Replace(line, " ", "", -1), ",")
		store = append(store, t[1])
	}
	if len(store)%3 != 0 {
		return errors.New(fmt.Sprintf("Store.csvが不正です(要素の数が合いません) : %s", src))
	}

	// 3行を1行にする
	var result []string
	for i := 0; i < len(store); i += 3 {
		result = append(result, fmt.Sprintf("%s %s %s", store[i], store[i+1], store[i+2]))
	}

	// 連結
	text := strings.Join(result, "\n")

	if err := ioutil.WriteFile(dst, []byte(text), 0644); err != nil {
		return err
	}
	return nil
}

func getACEScript(s string, r Range) []byte {
	return []byte(fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
	set www [ sx_find_wave_in_file $xml %s ]
	sx_export_csv on
	sx_export_range %s %s %s
	sx_export_data  "store.csv" $www
	`, s, r.Start, r.Stop, r.Step))
}

func getSPIScript(s Simulation, monte string, addfile string) ([]byte, error) {
	p := filepath.Join(ConfigDir, "spitemplate.spi")
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return []byte{}, err
	}
	tmplt := string(b)
	return []byte(fmt.Sprintf(tmplt, s.Vtn.Voltage, s.Vtn.Sigma, s.Vtn.Deviation,
		s.Vtp.Voltage, s.Vtp.Sigma, s.Vtp.Deviation, addfile, monte)), nil
}

func tryMkdir(p string) error {
	if _, err := os.Stat(p); err != nil {
		if e := os.MkdirAll(p, 0755); e != nil {
			return e
		}
		log.Print("Mkdir : ", p)
	}
	return nil
}

func sendTaskFile(src string, status bool) error {
	base := filepath.Base(src)
	dst := filepath.Join(DoneRunDir, base)
	if status {
		dst = filepath.Join(FailedRunDir, base)
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	log.Println("Move to ", dst)
	return nil
}

func moveTo(from string, f string, dir string) {
	src := filepath.Join(from, f)

	dst := filepath.Join(dir, f)

	if err := os.Rename(src, dst); err != nil {
		log.Fatal(err)
	}

	log.Print("Move to ", dst)
}

var ContinueWhenFaild, GargabeCollect bool

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().IntP("number", "n", 1, "実行するシミュレーションセットの個数です")
	runCmd.PersistentFlags().BoolVarP(&ContinueWhenFaild, "continue", "C", false, "連続して実行する時、どれかがコケても次のシミュレーションを行います")
	runCmd.PersistentFlags().Bool("all", false, "全部実行します")
	runCmd.PersistentFlags().BoolVar(&SlackNoNotify, "no-notify", false, "Slackに通知しません")
	runCmd.PersistentFlags().BoolVar(&GargabeCollect, "GC", false, "結果のCSV以外を削除します")
}
