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
	"fmt"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("make called")
	},
}

func makeTask() Task {
	t := Task{
		Config: config,
	}

	//Vtp
	t.Config.Vtp.Voltage = strconv.ParseFloat(getValue("Vtpの中央値です(default : ", t.Config.Vtp.Voltage, ")", string(t.Config.Vtp.Voltage)), 64)
	t.Config.Vtp.Sigma = strconv.ParseFloat(getValue("Vtpの中央値です(default : ", t.Config.Vtp.Sigma, ")", string(t.Config.Vtp.Sigma)), 64)
	t.Config.Vtp.Deviation = strconv.ParseFloat(getValue("Vtpの中央値です(default : ", t.Config.Vtp.Deviation, ")", string(t.Config.Vtp.Deviation)), 64)

	//Vtn
	t.Config.Vtn.Voltage = strconv.ParseFloat(getValue("Vtnの中央値です(default : ", t.Config.Vtn.Voltage, ")", string(t.Config.Vtn.Voltage)), 64)
	t.Config.Vtn.Sigma = strconv.ParseFloat(getValue("Vtnの中央値です(default : ", t.Config.Vtn.Sigma, ")", string(t.Config.Vtn.Sigma)), 64)
	t.Config.Vtn.Deviation = strconv.ParseFloat(getValue("Vtnの中央値です(default : ", t.Config.Vtn.Deviation, ")", string(t.Config.Vtn.Deviation)), 64)

	//Monte
	fmt.Printf("モンテカルロの回数をカンマ区切りで入力してください(default : %v)\n", t.Config.Monte)
	if res := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task")); len(res) != 0 {
		t.Config.Monte = strings.Split(res, ",")
	}

	//Range
	t.Config.Range.Start = getValue("書き出しを開始する時間です(")

}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasSuffix(s, in.GetWordBeforeCursor(), true)
}

func getValue(ask string, def string) string {
	fmt.Println(ask)

	res := prompt.Input(">>> ", completer, prompt.OptionTitle("UHA make Task"))
	if len(res) == 0 {
		return def
	}
	return res
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
