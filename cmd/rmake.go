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
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mattn/go-pipeline"
	"github.com/spf13/cobra"
)

// rmakeCmd represents the rmake command
var rmakeCmd = &cobra.Command{
	Use:   "rmake",
	Short: "sigmaの範囲を指定してタスクを連続生成",
	Long: `sigmaの範囲を指定してタスクを連続生成します。
	
Usage:
	UHA rmake start step stop
	
これ以外の値はデフォルト値が使われます`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("requires at less 3 args [start,step,stop]")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var num bool
		num, err = cmd.PersistentFlags().GetBool("number")
		if err != nil {
			log.Fatal(err)
		}

		var steps []string
		var start float64
		var step float64

		start, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		step, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		if num {
			cnt, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			steps = makeStepN(start, step, cnt)
		} else {
			stop, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			steps = makeStep(start, step, stop)
		}

		if err := rmakeTask(steps); err != nil {
			log.Fatal(err)
		}

	},
}

func makeStepN(start float64, step float64, n int64) []string {
	b, err := pipeline.Output(
		[]string{"seq", fmt.Sprint(start), fmt.Sprint(step), fmt.Sprint(1)},
		[]string{"head", "-n", fmt.Sprint(n)},
	)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(b), "\n")
}

func makeStep(start float64, step float64, stop float64) []string {
	b, err := pipeline.Output(
		[]string{"seq", fmt.Sprint(start), fmt.Sprint(step), fmt.Sprint(stop)},
	)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(b), "\n")
}

func rmakeTask(list []string) error {
	for _, v := range list {
		if len(v) == 0 {
			continue
		}
		sigma, _ := strconv.ParseFloat(v, 64)
		t := Task{
			Simulation: config.Simulation,
		}

		t.Simulation.Vtn.Sigma = sigma
		t.Simulation.Vtp.Sigma = sigma

		if err := writeTask(t); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(rmakeCmd)
	rmakeCmd.PersistentFlags().BoolP("number", "n", false, "始点、刻み幅、回数を指定して生成します")
}
