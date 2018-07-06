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
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "タスクを生成します",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		makeTask()
	},
}

func makeTask() {
	t := Task{
		Config: config,
	}

	//Vtp
	pv := getValue(fmt.Sprintf("Vtpの電圧です(default : %.4f)\n", t.Config.Vtp.Voltage), fmt.Sprint(t.Config.Vtp.Voltage))
	var pe error
	t.Config.Vtp.Voltage, pe = strconv.ParseFloat(pv, 64)
	if pe != nil {
		log.Fatal(pe)
	}
	ps := getValue(fmt.Sprintf("Vtpのシグマです(default : %.4f)\n", t.Config.Vtp.Sigma), fmt.Sprint(t.Config.Vtp.Sigma))
	t.Config.Vtp.Sigma, pe = strconv.ParseFloat(ps, 64)
	if pe != nil {
		log.Fatal(pe)
	}

	pd := getValue(fmt.Sprintf("Vtpの中央値です(default : %.4f)\n", t.Config.Vtp.Deviation), fmt.Sprint(t.Config.Vtp.Deviation))
	t.Config.Vtp.Deviation, pe = strconv.ParseFloat(pd, 64)
	if pe != nil {
		log.Fatal(pe)
	}
	//Vtn
	nv := getValue(fmt.Sprintf("Vtnの中央値です(default : %.4f)\n", t.Config.Vtn.Voltage), fmt.Sprint(t.Config.Vtn.Voltage))
	var ne error
	t.Config.Vtn.Voltage, ne = strconv.ParseFloat(nv, 64)
	if ne != nil {
		log.Fatal(ne)
	}
	ns := getValue(fmt.Sprintf("Vtnの中央値です(default : %.4f)\n", t.Config.Vtn.Sigma), fmt.Sprint(t.Config.Vtn.Sigma))
	t.Config.Vtn.Sigma, ne = strconv.ParseFloat(ns, 64)
	if ne != nil {
		log.Fatal(ne)
	}

	nd := getValue(fmt.Sprintf("Vtnの中央値です(default : %.4f)\n", t.Config.Vtn.Deviation), fmt.Sprint(t.Config.Vtn.Deviation))
	t.Config.Vtn.Deviation, ne = strconv.ParseFloat(nd, 64)
	if ne != nil {
		log.Fatal(ne)
	}

	//Monte
	fmt.Printf("モンテカルロの回数をカンマ区切りで入力してください(default : %v)\n", t.Config.Monte)
	if res := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task")); len(res) != 0 {
		t.Config.Monte = strings.Split(res, ",")
	}

	//Range
	t.Config.Range.Start = getValue(fmt.Sprintf("書き出しを開始する時間です(default %s)\n", t.Config.Range.Start), t.Config.Range.Start)
	t.Config.Range.Stop = getValue(fmt.Sprintf("書き出しを終了する時間です(default %s)\n", t.Config.Range.Stop), t.Config.Range.Stop)
	t.Config.Range.Step = getValue(fmt.Sprintf("書き出しの刻み幅です(default %s)\n", t.Config.Range.Step), t.Config.Range.Step)

	//DstDir
	t.Config.DstDir = getValue(fmt.Sprintf("結果が書き出されるディレクトリです(default %s)\n", t.Config.DstDir), t.Config.DstDir)
	//SimDir
	t.Config.SimDir = getValue(fmt.Sprintf("netlistが置かれるべきディレクトリです(default %s)\n", t.Config.SimDir), t.Config.SimDir)

	fmt.Printf("Vtn:AGAUSS(%.4f,%.4f,%.4f)\nVtn:AGAUSS(%.4f,%.4f,%.4f)\n", t.Config.Vtn.Voltage, t.Config.Vtn.Sigma, t.Config.Vtn.Deviation, t.Config.Vtn.Voltage, t.Config.Vtn.Sigma, t.Config.Vtn.Deviation)
	fmt.Printf("Monte:%v\n", t.Config.Monte)
	fmt.Printf("Range:[Start,Stop,Step] : %v\nDstDir:%s\nSimDir:%s\n", t.Config.Range, t.Config.DstDir, t.Config.SimDir)

	fmt.Println("この設定でシミュレーションタスクを発行します(y/n)")
	ans := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task Confirm"))
	if ans != "Y" && ans != "y" {
		log.Fatal("UHA make Task has canceled")
	}

	p := config.TaskDir
	if _, err := os.Stat(p); err != nil {
		log.Println("Can not find task dir : ", config.TaskDir)
		log.Println("Try to mkdir...")
		if e := os.Mkdir(config.TaskDir, 0755); e != nil {
			log.Fatal(e)
		}
	}

	j := path.Join(p, fmt.Sprint(time.Now().Format("20060102150405"), "_sigma", t.Config.Vtn.Sigma, ".json"))
	//f, err := os.OpenFile(j,os.O_CREATE|os.O_WRONLY,0644)
	b, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(j, b, os.ModePerm); err != nil {
		log.Fatal(err)
	}

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
}
