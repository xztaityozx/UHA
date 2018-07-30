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
	"time"

	"github.com/spf13/cobra"
)

// smakeCmd represents the smake command
var smakeCmd = &cobra.Command{
	Use:   "smake",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var times, number int
		var sigma float64
		var err error
		times, err = cmd.PersistentFlags().GetInt("times")
		if err != nil {
			log.Fatal(err)
		}
		number, err = cmd.PersistentFlags().GetInt("number")
		if err != nil {
			log.Fatal(err)
		}

		nt := NSeedTask{
			Simulation: config.Simulation,
			Count:      number,
		}

		sigma, err = cmd.PersistentFlags().GetFloat64("sigma")
		if err != nil {
			log.Fatal(err)
		}

		nt.Simulation.Monte = []string{strconv.Itoa(times)}
		nt.Simulation.Vtn.Sigma = sigma
		nt.Simulation.Vtp.Sigma = sigma

		if err := smakeTask(nt); err != nil {
			log.Fatal(err)
		}
	},
}

func smakeTask(nt NSeedTask) error {
	t := time.Now().Format("20060102150405")
	f := fmt.Sprintf("%s_N%d.json", t, nt.Count)

	p := filepath.Join(ReserveSRunDir, f)
	b, err := json.MarshalIndent(nt, "", "    ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(p, b, 0644); err != nil {
		return err
	}

	log.Println("Write SEED Task to : ", p)
	return nil
}

type NSeedTask struct {
	Simulation Simulation
	Count      int
}

func init() {
	rootCmd.AddCommand(smakeCmd)
	smakeCmd.PersistentFlags().IntP("times", "T", 50000, "モンテカルロの回数です")
	smakeCmd.PersistentFlags().Float64P("sigma", "S", config.Simulation.Vtn.Sigma, "シグマの値を指定します")
	smakeCmd.PersistentFlags().IntP("number", "n", 1, "SEEDの個数です")
}
