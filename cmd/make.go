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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "タスクを生成します",
	Long: `対話形式でシミュレーションセットを作成します。
作成されたセットは "UHA run"コマンドで実行することができます。

Usage : UHA make`,
	Run: func(cmd *cobra.Command, args []string) {
		sim := config.Simulation

		DstDir, _ := cmd.PersistentFlags().GetString("out")
		if len(DstDir) != 0 {
			sim.DstDir, _ = homedir.Expand(DstDir)
		}

		skip, _ := cmd.PersistentFlags().GetBool("default")
		yes, _ := cmd.PersistentFlags().GetBool("yes")

		if !skip {
			interactive(&sim, yes)
		}

		f, err := writeTask(sim)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Write Task File to : ", f)

	},
}

func interactive(sim *Simulation, yes bool) {
	// Vtp
	sim.Vtp.Voltage = getFloatVar("Vtpのしきい値電圧です", sim.Vtp.Voltage, "Vtp Volt")
	sim.Vtp.Sigma = getFloatVar("Vtpのシグマです", sim.Vtp.Sigma, "Vtp Sigma")
	sim.Vtp.Deviation = getFloatVar("Vtpの中央値です", sim.Vtp.Deviation, "Vtp Deviation")
	// Vtn
	sim.Vtn.Voltage = getFloatVar("Vtnのしきい値電圧です", sim.Vtn.Voltage, "Vtn Volt")
	sim.Vtn.Sigma = getFloatVar("Vtnのシグマです", sim.Vtn.Sigma, "Vtn Sigma")
	sim.Vtn.Deviation = getFloatVar("Vtnの中央値です", sim.Vtn.Deviation, "Vtn Deviation")
	// Monte
	sim.Monte = getStringSliceVar("モンテカルロの回数をカンマ区切りで入力します", sim.Monte, "Monte")
	// Range
	sim.Range.Start = getStringVar("プロットの開始時間です", sim.Range.Start, "Range Start")
	sim.Range.Stop = getStringVar("プロットの終了時間です", sim.Range.Stop, "Range Stop")
	sim.Range.Step = getStringVar("プロットの刻み幅です", sim.Range.Step, "Range Step")
	// Signal
	sim.Signal = getStringVar("プロットする信号線名です", sim.Signal, "Signal")
	// Dst
	sim.DstDir = getStringVar("結果が書き出される親ディレクトリです", sim.DstDir, "DstDir")
	// Sim
	sim.SimDir = getStringVar("netlistがあるディレクトリです", sim.SimDir, "SimDir")
	// SEED
	sim.SEED = getIntVar("SEED値です", sim.SEED, "SEED")

	fmt.Println("Vtp : ", sim.Vtp)
	fmt.Println("Vtn : ", sim.Vtn)
	fmt.Println("Range : ", sim.Range)
	fmt.Println("Monte : ", sim.Monte)
	fmt.Println("SEED : ", sim.SEED)
	fmt.Println("DstDir : ", sim.DstDir)
	fmt.Println("SimDir : ", sim.SimDir)

	if !yes {
		res := getStringVar("これでいいですか？", "no", "confirm")
		if res != "yes" {
			log.Fatal("中止しました")
		}
	}
}

func writeTask(sim Simulation) (string, error) {
	t := time.Now().Format("20060102150405")
	detail := fmt.Sprintf("VtpVolt%.4f_VtnVolt%.4f_Sigma%.4f_Monte%s_%s.json",
		sim.Vtp.Voltage,
		sim.Vtn.Voltage,
		sim.Vtn.Sigma,
		sim.Monte[0],
		sim.Monte[len(sim.Monte)-1])

	f := filepath.Join(ReserveRunDir, fmt.Sprintf("%s%s", t, detail))

	b, err := json.MarshalIndent(sim, "", "    ")
	if err != nil {
		return "", err
	}

	return f, ioutil.WriteFile(f, b, 0644)
}

func getStringVar(description string, def string, title string) string {
	res := prompt.Input(fmt.Sprintf("%s(default : %s)\n>>> ", description, def), completer, prompt.OptionTitle(fmt.Sprintf("UHA make %s", title)))
	if len(res) != 0 {
		return res
	}
	return def
}

func getIntVar(description string, def int, title string) int {
	res, err := strconv.Atoi(getStringVar(description, string(def), title))
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func getFloatVar(description string, def float64, title string) float64 {
	res, err := strconv.ParseFloat(getStringVar(description, fmt.Sprint(def), title), 64)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func getStringSliceVar(description string, def []string, title string) []string {
	return strings.Split(getStringVar(description, strings.Join(def, ","), title), ",")
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasSuffix(s, in.GetWordBeforeCursor(), true)
}

func init() {
	rootCmd.AddCommand(makeCmd)

	makeCmd.PersistentFlags().Float64("sigma", 0, "Vtp,VtnのSigmaを設定します")
	makeCmd.PersistentFlags().BoolP("default", "D", false, "設定ファイルをそのままタスクにします。オプションで値をしているとそちらが優先されます")
	makeCmd.PersistentFlags().BoolP("yes", "y", false, "y/nをスキップします")
	makeCmd.PersistentFlags().Float64P("VtpVoltage", "P", 0, "Vtpのしきい値電圧です")
	makeCmd.PersistentFlags().Float64P("VtnVoltage", "N", 0, "Vtnのしきい値電圧です")
	makeCmd.PersistentFlags().Int("SEED", 1, "SEED値です")
	makeCmd.PersistentFlags().StringP("out", "o", "", "書き出し先です")

	viper.BindPFlag("Simulation.Vtp.Voltage", makeCmd.PersistentFlags().Lookup("VtpVoltage"))
	viper.BindPFlag("Simulation.Vtn.Voltage", makeCmd.PersistentFlags().Lookup("VtnVoltage"))
	viper.BindPFlag("Simulation.Vtn.Sigma", makeCmd.PersistentFlags().Lookup("sigma"))
	viper.BindPFlag("Simulation.Vtp.Sigma", makeCmd.PersistentFlags().Lookup("sigma"))
	viper.BindPFlag("Simulation.SEED", makeCmd.PersistentFlags().Lookup("SEED"))
}
