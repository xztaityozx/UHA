// Copyright © 2019 xztaityozx
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
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// retryCmd represents the retry command
var retryCmd = &cobra.Command{
	Use:   "retry",
	Short: "失敗したシミュレーションを再実行するサブコマンドです",
	Long: `UHA retry [-N, --VtnVoltage [float]] [-P, --VtpVoltage [float]] [-C, --continue] [-p, --parallel [int]] [-T, --times [int]] [-S, --sigma [float]] [SEEDList file]

[SEEDList file]に記述されているSEEDを使ってシミュレーションを並列実行します
SEEDListのフォーマットはSEED値が1行1個書かれているテキストファイルです
省略するとStdInから読み取ります

ここで設定できない値はconfigの値と同じになります
`,
	Run: func(cmd *cobra.Command, args []string) {
		var gc, conti bool
		var err error
		var prlel, times int
		var vtn, vtp, sigma float64

		times, err = cmd.Flags().GetInt("times")
		if err != nil {
			log.Fatal(err)
		}
		vtp, err = cmd.Flags().GetFloat64("VtpVoltage")
		if err != nil {
			log.Fatal(err)
		}
		sigma, err = cmd.Flags().GetFloat64("sigma")
		if err != nil {
			log.Fatal(err)
		}
		vtn, err = cmd.Flags().GetFloat64("VtnVoltage")
		if err != nil {
			log.Fatal(err)
		}

		conti, err = cmd.Flags().GetBool("continue")
		if err != nil {
			log.Fatal(err)
		}
		prlel, err = cmd.Flags().GetInt("parallel")
		if err != nil {
			log.Fatal(err)
		}
		gc, err = cmd.Flags().GetBool("GC")
		if err != nil {
			log.Fatal(err)
		}

		file := os.Stdin
		if len(args) != 0 {
			file, err = os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}

		rt := NewRetryTask(file, times, prlel, vtn, vtp, sigma, gc, conti)
		msg := rt.Run()
		Post(config.SlackConfig, msg)
	},
}

type RetryTask struct {
	SEEDList   []int
	VtnVoltage float64
	VtpVoltage float64
	Sigma      float64
	Times      int
	Parallel   int
	GC         bool
	Continue   bool
}

func NewRetryTask(r io.Reader, t, p int, vtn, vtp, sigma float64, gc, con bool) RetryTask {

	rt := RetryTask{
		VtnVoltage: vtn,
		VtpVoltage: vtp,
		Sigma:      sigma,
		Times:      t,
		Parallel:   p,
		GC:         gc,
		Continue:   con,
	}

	// SEEDListを取り出す
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if p, err := strconv.Atoi(line); err != nil {
			log.Print("UHA retry:", line)
		} else {
			rt.SEEDList = append(rt.SEEDList, p)
		}
	}

	return rt
}

func (rt RetryTask) Run() SlackMessage {
	tasks := rt.BuildRangeSEEDTaskList()
	spin := spinner.New(spinner.CharSets[18], 50*time.Millisecond)
	spin.Suffix = " UHA retrying... "
	spin.FinalMSG = "All Retry had Finished\n"
	spin.Writer = os.Stderr
	spin.Start()
	defer spin.Stop()

	var wg sync.WaitGroup
	limit := make(chan struct{}, rt.Parallel)

	var summary []RetrySummary

	act := func(rst RangeSEEDTask, idx int) {
		limit <- struct{}{}
		defer wg.Done()

		log.Println("UHA retry: Start Task", idx)

		_, err := srun(rst, rt.GC)
		if err != nil {
			if rt.Continue {
				summary = append(summary, RetrySummary{rst.SEED, false})
				log.Println("Task", idx, "had failed...")
				log.Println(err)
			} else {
				PostFailed(config.SlackConfig, err)
				log.Fatal(err)
			}
		} else {
			summary = append(summary, RetrySummary{rst.SEED, true})
		}

		log.Printf("Finished (%d/%d)\n", idx, len(tasks))
		<-limit
	}

	for i, rst := range tasks {
		wg.Add(1)
		go act(rst, i+1)
	}

	wg.Wait()

	msg := SlackMessage{}

	for _, v := range summary {
		status := "Success"
		if !v.Status {
			status = "Failed"
		}
		fmt.Printf("%d: %s\n", v.SEED, status)
		if status == "Success" {
			msg.Succsess++
		} else {
			msg.Failed++
		}
	}

	return msg
}

type RetrySummary struct {
	SEED   int
	Status bool
}

func (rt RetryTask) BuildRangeSEEDTaskList() []RangeSEEDTask {
	var rst []RangeSEEDTask

	ace := getACEScript(config.Simulation.Signal, config.Simulation.Range)
	if err := ioutil.WriteFile(filepath.Join(config.Simulation.DstDir, "extract.ace"), ace, 0644); err != nil {
		log.Fatal(err)
	}

	for _, i := range rt.SEEDList {
		rst = append(rst, RangeSEEDTask{
			Monte:   fmt.Sprint(rt.Times),
			BaseDir: config.Simulation.DstDir,
			Sim:     config.Simulation.SimDir,
			SEED:    i,
			Vtn: Node{
				Voltage:   rt.VtnVoltage,
				Sigma:     rt.Sigma,
				Deviation: config.Simulation.Vtn.Deviation,
			},
			Vtp: Node{
				Voltage:   rt.VtpVoltage,
				Sigma:     rt.Sigma,
				Deviation: config.Simulation.Vtp.Deviation,
			},
			Sigma: rt.Sigma,
		})
	}

	return rst
}

// listUp サブコマンドP
var listupCmd = &cobra.Command{
	Use:   "list",
	Short: "失敗したSEEDをリストアップします",
	Long: `UHA retry listup [start int] [end int]

失敗したSEEDをリストアップします．
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("引数は2つ必要です")
		}

		start, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		end, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		files, err := ioutil.ReadDir(wd)
		dic := make(map[int]bool)
		for _, v := range files {
			name := v.Name()
			trimed := strings.Replace(name, "SEED", "", -1)
			trimed = strings.Replace(trimed, ".csv", "", -1)

			if seed, err := strconv.Atoi(trimed); err == nil {
				dic[seed] = true
			}
		}

		for i := start; i <= end; i++ {
			if !dic[i] {
				fmt.Println(i)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(retryCmd)
	retryCmd.AddCommand(listupCmd)

	retryCmd.Flags().Bool("GC", false, "最後に掃除をします")
	retryCmd.Flags().BoolP("continue", "C", false, "どこかでシミュレーションが失敗しても続けます")
	retryCmd.Flags().IntP("parallel", "p", 2, "並列実行する個数です")
	retryCmd.Flags().Float64P("VtnVoltage", "N", 0.6, "Vtnのしきい値電圧です")
	retryCmd.Flags().Float64P("VtpVoltage", "P", -0.6, "Vtpのしきい値電圧です")
	retryCmd.Flags().IntP("times", "T", 5000, "一回当たりのシミュレーション回数です")
	retryCmd.Flags().Float64P("sigma", "S", 0.046, "シグマの値です")
}
