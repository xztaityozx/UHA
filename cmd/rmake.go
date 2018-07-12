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
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// rmakeCmd represents the rmake command
var rmakeCmd = &cobra.Command{
	Use:   "rmake",
	Short: "",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("requires at less 3 args [start,step,stop]")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var start, stop, step float64
		var err error
		start, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		step, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		stop, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		if err := rmakeTask(start, step, stop); err != nil {
			log.Fatal(err)
		}
	},
}

func rmakeTask(start float64, step float64, stop float64) error {
	for ; start <= stop; start += step {
		t := Task{
			Simulation: config.Simulation,
		}
		t.Simulation.Vtn.Sigma = start
		t.Simulation.Vtp.Sigma = start

		if err := writeTask(t); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(rmakeCmd)
}
