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
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "タスクを生成します",
	Long: `対話形式でシミュレーションセットを作成します。
作成されたセットは "UHA run"コマンドで実行することができます。

Usage : UHA make`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		skip, err = cmd.PersistentFlags().GetBool("default")
		if err != nil {
			log.Fatal(err)
		}

		Sigma, err = cmd.PersistentFlags().GetFloat64("sigma")
		if err != nil {
			log.Fatal(err)
		}

		yes, err = cmd.PersistentFlags().GetBool("yes")
		if err != nil {
			log.Fatal(err)
		}
		makeTask()
	},
}

var Sigma float64
var skip bool
var yes bool

func interactive(t Task) Task {
	//Vtp
	pv := getValue(fmt.Sprintf("Vtpのしきい値電圧です(default : %.4f)\n", t.Simulation.Vtp.Voltage), fmt.Sprint(t.Simulation.Vtp.Voltage))
	var pe error
	t.Simulation.Vtp.Voltage, pe = strconv.ParseFloat(pv, 64)
	if pe != nil {
		log.Fatal(pe)
	}
	ps := getValue(fmt.Sprintf("Vtpのシグマです(default : %.4f)\n", t.Simulation.Vtp.Sigma), fmt.Sprint(t.Simulation.Vtp.Sigma))
	t.Simulation.Vtp.Sigma, pe = strconv.ParseFloat(ps, 64)
	if pe != nil {
		log.Fatal(pe)
	}

	pd := getValue(fmt.Sprintf("Vtpの中央値です(default : %.4f)\n", t.Simulation.Vtp.Deviation), fmt.Sprint(t.Simulation.Vtp.Deviation))
	t.Simulation.Vtp.Deviation, pe = strconv.ParseFloat(pd, 64)
	if pe != nil {
		log.Fatal(pe)
	}
	//Vtn
	nv := getValue(fmt.Sprintf("Vtnのしきい値電圧です(default : %.4f)\n", t.Simulation.Vtn.Voltage), fmt.Sprint(t.Simulation.Vtn.Voltage))
	var ne error
	t.Simulation.Vtn.Voltage, ne = strconv.ParseFloat(nv, 64)
	if ne != nil {
		log.Fatal(ne)
	}
	ns := getValue(fmt.Sprintf("Vtnのシグマです(default : %.4f)\n", t.Simulation.Vtn.Sigma), fmt.Sprint(t.Simulation.Vtn.Sigma))
	t.Simulation.Vtn.Sigma, ne = strconv.ParseFloat(ns, 64)
	if ne != nil {
		log.Fatal(ne)
	}

	nd := getValue(fmt.Sprintf("Vtnの中央値です(default : %.4f)\n", t.Simulation.Vtn.Deviation), fmt.Sprint(t.Simulation.Vtn.Deviation))
	t.Simulation.Vtn.Deviation, ne = strconv.ParseFloat(nd, 64)
	if ne != nil {
		log.Fatal(ne)
	}

	//Monte
	fmt.Printf("モンテカルロの回数をカンマ区切りで入力してください(default : %v)\n", t.Simulation.Monte)
	if res := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task")); len(res) != 0 {
		t.Simulation.Monte = strings.Split(res, ",")
	}

	//Range
	t.Simulation.Range.Start = getValue(fmt.Sprintf("書き出しを開始する時間です(default %s)\n", t.Simulation.Range.Start), t.Simulation.Range.Start)
	t.Simulation.Range.Stop = getValue(fmt.Sprintf("書き出しを終了する時間です(default %s)\n", t.Simulation.Range.Stop), t.Simulation.Range.Stop)
	t.Simulation.Range.Step = getValue(fmt.Sprintf("書き出しの刻み幅です(default %s)\n", t.Simulation.Range.Step), t.Simulation.Range.Step)

	//DstDir
	t.Simulation.DstDir = getValue(fmt.Sprintf("結果が書き出されるディレクトリです(default %s)\n", t.Simulation.DstDir), t.Simulation.DstDir)
	//SimDir
	t.Simulation.SimDir = getValue(fmt.Sprintf("netlistが置かれるべきディレクトリです(default %s)\n", t.Simulation.SimDir), t.Simulation.SimDir)

	//LibDir
	//t.Simulation.LibDir = getValue(fmt.Sprintf("ライブラリのディレクトリです(default %s)\n", t.Simulation.LibDir), t.Simulation.LibDir)
	//AddFile
	//t.Simulation.AddFile = getValue(fmt.Sprintf("addfileへのパスです(default %s)\n", t.Simulation.AddFile), t.Simulation.AddFile)
	//ModelFile
	//t.Simulation.ModelFile = getValue(fmt.Sprintf("modelfileへのパスです(default %s)\n", t.Simulation.ModelFile), t.Simulation.ModelFile)

	//Signal
	t.Simulation.Signal = getValue(fmt.Sprintf("プロットしたい信号線名です(default %s)\n", t.Simulation.Signal), t.Simulation.Signal)

	//SEED
	seed := getValue(fmt.Sprintf("SEED値です(default : %d)\n", t.Simulation.SEED), string(t.Simulation.SEED))
	if len(seed) != 0 {
		t.Simulation.SEED, _ = strconv.Atoi(seed)
	}

	return t
}

func makeTask() {
	t := Task{
		Simulation: config.Simulation,
	}

	if Sigma != 0.0 {
		t.Simulation.Vtn.Sigma = Sigma
		t.Simulation.Vtp.Sigma = Sigma
	}

	if !skip {
		t = interactive(t)
	}

	if !yes {

		fmt.Printf("Vtn:AGAUSS(%.4f,%.4f,%.4f)\nVtp:AGAUSS(%.4f,%.4f,%.4f)\n", t.Simulation.Vtn.Voltage, t.Simulation.Vtn.Sigma, t.Simulation.Vtn.Deviation, t.Simulation.Vtp.Voltage, t.Simulation.Vtp.Sigma, t.Simulation.Vtp.Deviation)
		fmt.Printf("Monte:%v\n", t.Simulation.Monte)
		fmt.Printf("Range:[Start,Stop,Step] : %v\nDstDir:%s\nSimDir:%s\n", t.Simulation.Range, t.Simulation.DstDir, t.Simulation.SimDir)

		fmt.Println("この設定でシミュレーションタスクを発行します(y/n)")
		ans := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task Confirm"))
		if ans != "Y" && ans != "y" {
			log.Fatal("UHA make Task has canceled")
		}
	}

	if err := writeTask(t); err != nil {
		log.Fatal(err)
	}
}

func writeTask(t Task) error {
	tryMkdir(ReserveRunDir)

	// ~ resolve
	t.Simulation.DstDir, _ = homedir.Expand(t.Simulation.DstDir)
	t.Simulation.SimDir, _ = homedir.Expand(t.Simulation.SimDir)
	//t.Simulation.AddFile, _ = homedir.Expand(t.Simulation.AddFile)
	//t.Simulation.LibDir, _ = homedir.Expand(t.Simulation.LibDir)
	//t.Simulation.ModelFile, _ = homedir.Expand(t.Simulation.ModelFile)

	j := path.Join(ReserveRunDir, fmt.Sprint(time.Now().Format("20060102150405"), "_sigma", t.Simulation.Vtn.Sigma, ".json"))
	//f, err := os.OpenFile(j,os.O_CREATE|os.O_WRONLY,0644)
	b, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(j, b, 0644); err != nil {
		return err
	}
	log.Println("Write Task to", j)
	return nil
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasSuffix(s, in.GetWordBeforeCursor(), true)
}

func getValue(ask string, def string) string {
	fmt.Print(ask)

	res := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task"))
	if len(res) == 0 {
		return def
	}
	return res
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.PersistentFlags().Float64("sigma", 0.0, "Vtp,VtnのSigmaを設定します")
	makeCmd.PersistentFlags().BoolP("default", "D", false, "設定ファイルをそのままタスクにします。オプションで値をしているとそちらが優先されます")
	makeCmd.PersistentFlags().BoolP("yes", "y", false, "y/nをスキップします")
}
