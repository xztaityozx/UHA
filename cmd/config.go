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
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("UHA config:\nCurrent Config:\n\tMonte:%v\n\tRange:%v\n\tDstDir:%s\n\tSimDir:%s\n", config.Monte, config.Range, config.Dstdir, config.Dstdir)
		fmt.Println("設定を変更します。空白を入力すると、現在の値を使います")

		c := Config{
			Dstdir: config.Dstdir,
			SimDir: config.SimDir,
			Monte:  config.Monte,
			Range:  config.Range,
		}

		// Monte
		fmt.Println("モンテカルロの回数をカンマ区切りで指定してください")
		m := strings.Split(prompt.Input(">>> ", completer, prompt.OptionTitle("UHA config")), ",")

		if len(m) != 0 {
			c.Monte = m
		}

		//Range
		fmt.Println("データを書き出す時の開始時間")
		start := prompt.Input(">>> ", completer)
		fmt.Println("データを書き出す時の終了時間")
		stop := prompt.Input(">>> ", completer)
		fmt.Println("データを取る間隔")
		step := prompt.Input(">>> ", completer)

		if len(start) != 0 {
			c.Range.Start = start
		}
		if len(stop) != 0 {
			c.Range.Stop = stop
		}
		if len(step) != 0 {
			c.Range.Step = step
		}

		//Dstdir
		fmt.Println("データを書き出すディレクトリ")
		dstdir := prompt.Input(">>> ", completer)

		//SimDir
		fmt.Println("netlistがあるディレクトリ")
		simdir := prompt.Input(">>> ", completer)

		if len(dstdir) != 0 {
			c.Dstdir = dstdir
		}
		if len(simdir) != 0 {
			c.SimDir = simdir
		}

		SaveConfig(c)
	},
}

type Task struct {
	Vtn     string
	Vtp     string
	SigName string
	Config  Config
}

type Config struct {
	SimDir      string
	Dstdir      string
	Monte       []string
	Range       Range
	Repositorys []Repository
}

//func makeTask() Task, error {

//}

func init() {
	rootCmd.AddCommand(configCmd)
}
