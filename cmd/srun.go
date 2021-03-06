// Copyright © 2018 xztaityozx
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
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// srunCmd represents the srun command
var srunCmd = &cobra.Command{
	Use:   "srun",
	Short: "smakeで作ったタスクを実行します",
	Long: `SEEDを連番生成しながら複数回モンテカルロを実行します
	
Usage:
	UHA srun [--number,-n [NUM]|--parallel,-P [NUM]|--all|--custom [commands]|--continue,-C|--no-notify]

	先に"UHA smake"でタスクを作ってから実行してください
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var conti, all, summary, gc bool
		var prlel, num, start int
		var err error
		conti, err = cmd.PersistentFlags().GetBool("continue")
		if err != nil {
			log.Fatal(err)
		}
		prlel, err = cmd.PersistentFlags().GetInt("parallel")
		if err != nil {
			log.Fatal(err)
		}
		all, err = cmd.PersistentFlags().GetBool("all")
		if err != nil {
			log.Fatal(err)
		}
		summary, err = cmd.PersistentFlags().GetBool("summary")
		if err != nil {
			log.Fatal(err)
		}
		gc, err = cmd.PersistentFlags().GetBool("GC")
		if err != nil {
			log.Fatal(err)
		}
		num, err = cmd.PersistentFlags().GetInt("number")
		if err != nil {
			log.Fatal(err)
		}
		start, err = cmd.PersistentFlags().GetInt("start")
		if err != nil {
			log.Fatal(err)
		}

		msg := SlackMessage{
			StartTime: time.Now(),
		}
		res := RunRangeSEEDSimulation(start, prlel, conti, all, gc, num)
		msg.FinishedTime = time.Now()

		if summary {
			printSummary(&res)
		}

		msg.Succsess, msg.Failed = countSuccessesRangeSEED(&res)

		if err = Post(config.SlackConfig, msg); err != nil {
			log.Fatal(err)
		} else if !SlackNoNotify {
			log.Println("Post To Slack")
		}
	},
}

func countSuccessesRangeSEED(summary *[]SRunSummary) (int, int) {
	suc := 0
	for _, v := range *summary {
		if v.Status {
			suc++
		}
	}
	return suc, len(*summary) - suc
}

type RangeSEEDTask struct {
	Addfile string
	SPI     string
	Dst     string
	BaseDir string
	Sim     string
	SEED    int
	Sigma   float64
	Vtp     Node
	Vtn     Node
	Monte   string
}

type SRunSummary struct {
	Name       string
	Status     bool
	StartTime  time.Time
	FinishTime time.Time
}

func srun(task RangeSEEDTask, gc bool) (SRunSummary, error) {
	var summary SRunSummary = SRunSummary{
		Name:      fmt.Sprintf("Sigma%.4f-SEED%03d", task.Sigma, task.SEED),
		StartTime: time.Now(),
		Status:    false,
	}

	// ディレクトリを作る
	if err := tryMkRangeSEEDDstDir(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}
	// Addfile
	if err := writeRangeSEEAddfile(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}
	// SPI
	if err := writeRangeSEEDSPI(&task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	// XML
	if err := copyRangeSEEDXmls(task); err != nil {
		summary.FinishTime = time.Now()
		return summary, err
	}

	// ゴミ掃除
	if gc {
		defer removeRangeSEEDGarbage(task)
	}

	// シミュレーション
	command := makeRangeSEEDCommand(task)
	err := exec.Command("bash", "-c", command).Run()
	summary.FinishTime = time.Now()
	if err == nil {
		summary.Status = true
	}

	return summary, err
}

// SRun本体
func RunRangeSEEDSimulation(start int, prlel int, conti bool, all bool, gc bool, num int) []SRunSummary {
	// Summary
	var rt []SRunSummary

	files, err := ioutil.ReadDir(ReserveSRunDir)
	if err != nil {
		PostFailed(config.SlackConfig, err)
		log.Fatal(err)
	}

	// シミュレーションする個数
	length := num
	if len(files) < num || all {
		length = len(files)
	}

	// spinner
	spin := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	spin.Suffix = " Running... "
	spin.FinalMSG = "All Task had Finished\n"
	spin.Writer = os.Stderr
	spin.Start()
	defer spin.Stop()

	// waitgroup
	var wg sync.WaitGroup
	limit := make(chan struct{}, prlel)

	for i := 0; i < length; i++ {
		// タスク読み出し
		tasks, file, err := readRangeSEEDTask(start)
		if err != nil {
			// 失敗しても継続)
			if conti {
				continue
			}
			PostFailed(config.SlackConfig, err)
			log.Fatal(err)
		}
		// 正常終了したか？
		success := true
		// 並列化
		for k, v := range tasks {
			wg.Add(1)
			go func(num int, v RangeSEEDTask) {
				limit <- struct{}{}
				defer wg.Done()

				sum, err := srun(v, gc)
				if err != nil {
					// 失敗
					success = false
					if conti {
						log.Println("Task", num, "had failed...")
					} else {
						PostFailed(config.SlackConfig, err)
						log.Fatal(err)
					}
				}
				rt = append(rt, sum)

				log.Printf("Finished (%d/%d)\n", num, len(tasks))
				<-limit
			}(k+1, v)
		}

		wg.Wait()
		to := DoneSRunDir
		if !success {
			to = FailedSRunDir
		}
		// タスクファイルを移動
		moveTo(ReserveSRunDir, file, to)
	}

	return rt
}

func printSummary(summarys *[]SRunSummary) {
	status := map[bool]string{
		true:  "\033[1:32m●\033[1:39m  ",
		false: "\033[1:31m●\033[1:39m  ",
	}

	fmt.Println("Name\tStatus\tStartTime\tFinishTime")

	for _, v := range *summarys {
		fmt.Printf("%s\t%s\t%s\t%s\n", v.Name, status[v.Status], v.StartTime.Format("2006/01/02/15:04.05"), v.FinishTime.Format("2006/01/02/15:04.05"))
	}

}

func makeRangeSEEDCommand(rst RangeSEEDTask) string {
	str := fmt.Sprintf("cd %s && hspice -mt 2 -i %s -o ./hspice &> ./hspice.log &&", rst.Dst, rst.SPI)
	str += fmt.Sprintf("wv -k -ace_no_gui ../../extract.ace &> ./wv.log &&")
	str += fmt.Sprintf("cat store.csv | sed '/^#/d;1,1d' | awk -F, '{print $2}' | xargs -n3 > ../Result/SEED%03d.csv ", rst.SEED)

	return str
}

func removeRangeSEEDGarbage(rst RangeSEEDTask) error {
	// remove Addfile
	if err := os.Remove(rst.Addfile); err != nil {
		return err
	}
	// remove SPI
	if err := os.Remove(rst.SPI); err != nil {
		return err
	}
	// remove Dst
	if err := os.RemoveAll(rst.Dst); err != nil {
		return err
	}
	return nil
}

func readRangeSEEDTask(start int) ([]RangeSEEDTask, string, error) {
	var nt NSeedTask
	files, err := ioutil.ReadDir(ReserveSRunDir)
	if err != nil {
		return nil, "", err
	}

	if len(files) == 0 {
		return nil, "", errors.New("タスクがありません")
	}

	p := filepath.Join(ReserveSRunDir, files[0].Name())
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, "", err
	}

	if err := json.Unmarshal(b, &nt); err != nil {
		return nil, "", err
	}
	ace := getACEScript(nt.Simulation.Signal, nt.Simulation.Range)
	if err := ioutil.WriteFile(filepath.Join(nt.Simulation.DstDir, "extract.ace"), ace, 0644); err != nil {
		return nil, "", err
	}

	var rt []RangeSEEDTask

	for i := start; i < nt.Count+start; i++ {
		rst := RangeSEEDTask{
			Monte:   nt.Simulation.Monte[0],
			BaseDir: nt.Simulation.DstDir,
			Sim:     nt.Simulation.SimDir,
			SEED:    i,
			Vtn:     nt.Simulation.Vtn,
			Vtp:     nt.Simulation.Vtp,
			Sigma:   nt.Simulation.Vtn.Sigma,
		}
		rt = append(rt, rst)
	}

	return rt, files[0].Name(), nil

}

func writeRangeSEEDSPI(rst *RangeSEEDTask) error {
	f := filepath.Join(rst.Sim, fmt.Sprintf("Vtn%.4fVtp%.4fMonte%s_Sigma%.4f_SEED%03d.spi", rst.Vtn.Voltage, rst.Vtp.Voltage, rst.Monte, rst.Sigma, rst.SEED))

	// SPI文字列を作る
	p := filepath.Join(ConfigDir, "spitemplate.spi")
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	tmplt := string(b)
	spi := []byte(fmt.Sprintf(tmplt,
		rst.Vtn.Voltage,
		rst.Vtn.Sigma,
		rst.Vtn.Deviation,
		rst.Vtp.Voltage,
		rst.Vtp.Sigma,
		rst.Vtp.Deviation,
		rst.Addfile,
		rst.Monte))

	if err := ioutil.WriteFile(f, spi, 0644); err != nil {
		return err
	}
	rst.SPI = f
	//log.Println("Write SPIscript To :", f)
	return nil
}

// このシミュレーションで使うAddfileを作る
func writeRangeSEEAddfile(rst *RangeSEEDTask) error {
	// tryMkdir
	dir := filepath.Join(rst.BaseDir, "Addfiles")
	if err := tryMkdir(dir); err != nil {
		return err
	}

	// テンプレを読む
	p := filepath.Join(ConfigDir, "addfile.txt")
	tmp, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	// addfileの中身をつくる
	addfile := []byte(fmt.Sprintf(string(tmp), rst.SEED))
	f := filepath.Join(dir, fmt.Sprintf("addfile%03d.txt", rst.SEED))

	if err := ioutil.WriteFile(f, addfile, 0644); err != nil {
		return err
	}

	// rstに設定して終わる
	rst.Addfile = f
	//log.Println("Write Addfile To :", f)
	return nil
}

// このシミュレーションの結果を書き出すディレクトリを作る
func tryMkRangeSEEDDstDir(rst *RangeSEEDTask) error {
	ts := filepath.Join(rst.BaseDir, fmt.Sprintf("RangeSEED_Vtn%.4fVtp%.4f_Sigma%.4f_Monte%s", rst.Vtn.Voltage, rst.Vtp.Voltage, rst.Sigma, rst.Monte))
	p := filepath.Join(ts, fmt.Sprintf("SEED%03d", rst.SEED))
	if err := tryMkdirSuppress(p); err != nil {
		return err
	}

	result := filepath.Join(ts, "Result")
	if err := tryMkdirSuppress(result); err != nil {
		return err
	}

	rst.Dst = p
	return nil
}

func tryMkdirSuppress(p string) error {
	if _, err := os.Stat(p); err != nil {
		if e := os.MkdirAll(p, 0755); e != nil {
			return e
		}
	}
	return nil
}

func copyRangeSEEDXmls(rst RangeSEEDTask) error {
	resultsxml := filepath.Join(SelfPath, "templates", rst.Monte)
	mapxml := filepath.Join(SelfPath, "templates", "resultsMap.xml")

	{
		src, e1 := ioutil.ReadFile(resultsxml)
		if e1 != nil {
			return e1
		}
		dst := filepath.Join(rst.Dst, "results.xml")

		if err := ioutil.WriteFile(dst, src, 0644); err != nil {
			return err
		}
	}
	{
		src, e1 := ioutil.ReadFile(mapxml)
		if e1 != nil {
			return e1
		}
		dst := filepath.Join(rst.Dst, "resultsMap.xml")
		if err := ioutil.WriteFile(dst, src, 0644); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(srunCmd)
	srunCmd.PersistentFlags().Bool("all", false, "すべて実行します")
	srunCmd.PersistentFlags().BoolP("continue", "C", false, "どこかでシミュレーションが失敗しても続けます")
	srunCmd.PersistentFlags().IntP("number", "n", 1, "実行するタスクの個数です。default : 1")
	srunCmd.PersistentFlags().IntP("parallel", "P", 2, "並列実行する個数です。default : 2")
	srunCmd.PersistentFlags().Int("start", 1, "SEEDの最初の値です")
	srunCmd.PersistentFlags().BoolVar(&SlackNoNotify, "no-notify", false, "Slackに通知しません")
}
