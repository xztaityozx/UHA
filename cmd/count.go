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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-pipeline"
	"github.com/spakin/awk"
	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "データを数え上げます",
	Long: `現在のディレクトリにあるCSVを見つけて、不良数を数え上げます
Usage:
	UHA count
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		if len(args) != 0 {
			dir = args[0]
		}

		agg, err := cmd.PersistentFlags().GetBool("aggregate")
		if err != nil {
			log.Fatal(err)
		}
		only, err := cmd.PersistentFlags().GetBool("only")
		if err != nil {
			log.Fatal(err)
		}

		if agg {
			r, f, err := dirAggregate(dir)
			if err != nil {
				log.Fatal(err)
			}
			if only {
				fmt.Println(f)
			} else {
				fmt.Println(r, f)
			}
		} else {
			for _, v := range Count(dir) {
				fmt.Println(v)
			}
		}
	},
}

func countup(p string) (int, int, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return -1, -1, err
	}
	s := awk.NewScript()
	s.Begin = func(s *awk.Script) { s.State = 0 }
	s.AppendStmt(
		func(s *awk.Script) bool {
			return s.F(0).Float64() >= 0.4 && s.F(3).Float64() >= 0.4
		}, func(s *awk.Script) {
			s.State = s.State.(int) + 1
		})

	r := strings.NewReader(string(b))
	if err := s.Run(r); err != nil {
		return s.NR, -1, err
	}

	return s.NR, s.State.(int), nil

}

func Count(wd string) []interface{} {
	var rt []interface{}

	if err := os.Chdir(wd); err != nil {
		log.Fatal(err)
	}

	wl := len(wd)
	sigma := wd[wl-6 : wl-1]

	if RangeSEEDCount {
		c := exec.Command("bash", "-c", "cd ../ && basename $(pwd) | sed 's/RangeSEED\\|_\\|Sigma\\|Monte.*$//g'")
		o, err := c.CombinedOutput()
		if err != nil {
			log.Fatal(string(o))
		}

		sigma = string(o)
	}

	rt = append(rt, sigma)

	// 数え上げ
	b, err := pipeline.Output(
		[]string{"ls", "-1"},
		[]string{"grep", ".csv"},
		[]string{"sort", "-n"},
	)
	if err != nil {
		log.Fatal(err)
	}

	files := strings.Split(string(b), "\n")
	for _, v := range files {
		if len(v) == 0 {
			continue
		}
		_, cnt, err := countup(filepath.Join(wd, v))
		if err != nil {
			log.Fatal(err)
		}

		rt = append(rt, cnt)
		//if only {
		//fmt.Println(cnt)
		//} else {
		//fmt.Println(v, cnt)
		//}
	}
	return rt
}

// ディレクトリ
func dirAggregate(dir string) (int, int, error) {
	// 移動する
	if err := os.Chdir(dir); err != nil {
		return -1, -1, err
	}

	b, err := pipeline.Output(
		[]string{"ls", "-1"},
		[]string{"grep", ".csv"},
		[]string{"sort", "-n"},
	)
	if err != nil {
		log.Fatal(err)
	}

	files := strings.Split(string(b), "\n")

	size := 0
	failure := 0

	var wg sync.WaitGroup
	s := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	s.Suffix = " Counting..."
	s.FinalMSG = "Aggregated\n"
	s.Start()
	defer s.Stop()

	for _, v := range files {
		if len(v) == 0 {
			continue
		}
		wg.Add(1)

		go func(v string) {
			defer wg.Done()
			n, cnt, err := countup(filepath.Join(dir, v))
			if err != nil {
				log.Fatal(err)
			}
			size += n
			failure += cnt
		}(v)
	}
	wg.Wait()

	return size, failure, nil
}

var RangeSEEDCount bool

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.PersistentFlags().BoolP("aggregate", "A", false, "ディレクトリ以下のファイルを1つのデータの集合としてカウントします")
	countCmd.PersistentFlags().BoolP("only", "o", false, "不良数だけを出力します")
	countCmd.PersistentFlags().BoolVarP(&RangeSEEDCount, "RangeSEED", "R", false, "RangeSEEDシミュレーションの結果を数え上げます")
}
