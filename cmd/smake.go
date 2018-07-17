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

		nt.Simulation.Monte = []string{strconv.Itoa(times)}

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

func makeSAddfile(seed int) (string, error) {
	p := filepath.Join(ConfigDir, "addfile.txt")
	tmp, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(string(tmp), seed), nil
}

type NSeedTask struct {
	Simulation Simulation
	Count      int
}

func init() {
	rootCmd.AddCommand(smakeCmd)
	smakeCmd.PersistentFlags().IntP("times", "T", 50000, "モンテカルロの回数です (default : 50000)")
	smakeCmd.PersistentFlags().IntP("number", "n", 1, "SEEDを変更する回数です")
}
