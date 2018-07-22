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
	UHA srun [--number,-n [NUM]|--parallel,-P [NUM]|--all|--custom [commands]|--continue,-C]

	先に"UHA smake"でタスクを作ってから実行してください
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conti, _ := cmd.PersistentFlags().GetBool("continue")
		prlel, _ := cmd.PersistentFlags().GetInt("parallel")
		all, _ := cmd.PersistentFlags().GetBool("all")
		custom, _ := cmd.PersistentFlags().GetStringSlice("custom")
		num, _ := cmd.PersistentFlags().GetInt("number")

		// task list
		list, files := readNSTaskFileList()
		if len(list) == 0 && len(custom) == 0 {
			log.Fatal("タスクが見つかりませんでした")
		}

		if !all {
			list = list[0:num]
		}

		if err := srun(prlel, conti, list, custom); err != nil {
			for _, v := range files {
				p := filepath.Join(ReserveSRunDir, v)
				moveTo(p, FailedSRunDir)
			}
		} else {
			for _, v := range files {
				p := filepath.Join(ReserveSRunDir, v)
				moveTo(p, DoneSRunDir)
			}
		}

	},
}

// prlel個並列にタスクを実行する。
func srun(prlel int, conti bool, tasks []NSeedTask, Custom []string) error {
	var commands []string

	if len(Custom) == 0 {
		for _, v := range tasks {
			res := makeSRun(v)
			commands = append(commands, res...)
		}
	} else {
		commands = Custom
	}

	log.Println(commands)

	log.Println("Start Simulation Set :Range=", len(commands))

	// WaitGroup
	var wg sync.WaitGroup
	limit := make(chan struct{}, prlel)
	count := 0

	// resultDir
	for _, v := range tasks {
		resultDir := filepath.Join(v.Simulation.DstDir, fmt.Sprintf("Sigma%.4f", v.Simulation.Vtn.Sigma))
		if err := tryMkdir(resultDir); err != nil {
			return err
		}
	}

	// スピナー
	s := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	s.Suffix = "Running... "
	s.FinalMSG = "Finished!"
	s.Start()

	for _, command := range commands {
		wg.Add(1)
		count++

		log.Println(command)
		flag := false

		go func(command string, cnt int) {
			limit <- struct{}{}
			defer wg.Done()

			if err := exec.Command("bash", "-c", command).Run(); err != nil {
				if !conti {
					log.Fatal(err)
				}
				flag = true
			} else {
				log.Printf("Finished (%d/%d)\n", cnt, len(commands))
			}
			<-limit
		}(command, count)

		if flag {
			return errors.New("Failed Simulation")
		}
	}
	wg.Wait()
	s.Stop()

	return nil
}

// ReserveSRunDirから、NSeedTaskのJSONとして正しいやつだけ列挙
func readNSTaskFileList() ([]NSeedTask, []string) {
	var rt []NSeedTask
	files, err := ioutil.ReadDir(ReserveSRunDir)
	if err != nil {
		log.Fatal(err)
	}

	var list []string

	for _, f := range files {
		p := filepath.Join(ReserveSRunDir, f.Name())
		b, rerr := ioutil.ReadFile(p)
		if rerr != nil {
			log.Fatal(rerr)
		}
		var nt NSeedTask
		jerr := json.Unmarshal(b, &nt)
		if jerr != nil {
			log.Println(jerr)
			continue
		}

		rt = append(rt, nt)
		list = append(list, f.Name())
	}

	return rt, list
}

func setResultDir(nt NSeedTask) error {
	for i := 1; i <= nt.Count; i++ {
		p := filepath.Join(nt.Simulation.DstDir, fmt.Sprintf("Monte%s_SEED%d", nt.Simulation.Monte[0], i))
		if err := tryMkdir(p); err != nil {
			return err
		}
	}
	return nil
}

func makeSRun(nt NSeedTask) []string {
	var rt []string

	addfile := nt.Simulation.SimDir
	// ディレクトリを作る
	if err := setResultDir(nt); err != nil {
		log.Fatal(err)
	}
	// Addfileを作る
	if err := setAddfileTo(nt.Count, addfile); err != nil {
		log.Fatal(err)
	}
	// SPIをつくる
	if err := setSEEDInputSPI(nt.Count, nt.Simulation.SimDir, nt.Simulation); err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= nt.Count; i++ {
		dst := filepath.Join(nt.Simulation.DstDir, fmt.Sprintf("Monte%s_SEED%d", nt.Simulation.Monte[0], i))
		input := filepath.Join(nt.Simulation.SimDir, fmt.Sprintf("%s_SEED%d_input.spi", nt.Simulation.Monte[0], i))

		str := fmt.Sprintf("cd %s && hspice -hpp -mt 4 -i %s -o ./hspice &> ./hspice.log && wv -k -ace_no_gui ../extract.ace &> wv.log && ", dst, input)
		str += fmt.Sprintf("cat store.csv | sed '/^#/d;1,1d' | awk -F, '{print $2}' | xargs -n3 >> ../Sigma%.4f/result\n", nt.Simulation.Vtn.Sigma)

		rt = append(rt, str)
	}

	return rt
}

func setSEEDInputSPI(cnt int, p string, sim Simulation) error {

	for i := 1; i <= cnt; i++ {
		spi, err := getSPIScript(sim, sim.Monte[0], fmt.Sprintf("addfile%d.txt", i))
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(p, fmt.Sprintf("%s_SEED%d_input.spi", sim.Monte[0], i)), spi, 0644); err != nil {
			return err
		}
	}

	return nil
}

func setAddfileTo(cnt int, p string) error {
	for i := 1; i <= cnt; i++ {
		s, err := makeSAddfile(i)
		b := []byte(s)
		if err != nil {
			return err
		}

		f := filepath.Join(p, fmt.Sprintf("addfile%d.txt", i))
		if err := ioutil.WriteFile(f, b, 0644); err != nil {
			return err
		}
		log.Printf("Write Addfile%d To : %s\n", i, f)
	}
	return nil
}

// ConfigDir以下にあるaddfile.txtをテンプレートに、SEEDを変更したaddfileの文字列を作る
func makeSAddfile(seed int) (string, error) {
	p := filepath.Join(ConfigDir, "addfile.txt")
	tmp, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(string(tmp), seed), nil
}

func init() {
	rootCmd.AddCommand(srunCmd)
	srunCmd.PersistentFlags().Bool("all", false, "すべて実行します")
	srunCmd.PersistentFlags().BoolP("continue", "C", false, "どこかでシミュレーションが失敗しても続けます")
	srunCmd.PersistentFlags().IntP("number", "n", 1, "実行するタスクの個数です。default : 1")
	srunCmd.PersistentFlags().IntP("parallel", "P", 2, "並列実行する個数です。default : 2")
	srunCmd.PersistentFlags().StringSlice("custom", []string{}, "カスタムコマンドを並列実行します")
}
