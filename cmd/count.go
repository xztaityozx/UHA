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
	"path/filepath"
	"sort"
	"strconv"
	"sync"

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
		cum, _ := cmd.Flags().GetBool("CumulativeSum")
		fOnly, _ := cmd.Flags().GetBool("failure-only")

		wd, _ := os.Getwd()
		res := GetAggregateDataAll(wd)
		if cum {
			get := CumulativeSum(&res)
			PrintAggregateData(&get, fOnly)
			return
		}

		PrintAggregateData(&res, fOnly)

	},
}

type AggregateData struct {
	Lines    int
	Failure  int
	FileName string
}

func NewAggregateData(f string) (AggregateData, error) {
	fp, err := os.OpenFile(f, os.O_RDONLY, 0644)
	defer fp.Close()
	if err != nil {
		return AggregateData{}, err
	}

	s := awk.NewScript()
	s.Begin = func(s *awk.Script) {
		s.State = 0
	}
	s.AppendStmt(func(s *awk.Script) bool {
		return s.F(1).Float64() >= countFirstFilter &&
			s.F(2).Float64() >= countSecondFilter &&
			s.F(3).Float64() >= countThirdFilter
	}, func(s *awk.Script) {
		s.State = s.State.(int) + 1
	})
	if err := s.Run(fp); err != nil {
		return AggregateData{}, err
	}

	return AggregateData{
		Failure:  s.State.(int),
		Lines:    s.NR,
		FileName: filepath.Base(f),
	}, nil
}

// dirにあるcsvをAggregateDataにしてからファイル名でソートして返す
func GetAggregateDataAll(dir string) []AggregateData {
	var rt []AggregateData
	rec, fin := aggWorker(dir)
	for {
		select {
		case res := <-rec:
			rt = append(rt, res)
		case <-fin:
			sort.Slice(rt, func(i, j int) bool {
				return rt[i].FileName < rt[j].FileName
			})
			return rt
		}
	}
}

func GetSigma(dir string) float64 {
	if RangeSEEDCount {
		dir = filepath.Dir(dir)
		base := filepath.Base(dir)
		if len(base) < 40 {
			log.Fatal("Invaild dir name")
		}
		f, fe := strconv.ParseFloat(base[34:40], 64)
		if fe != nil {
			log.Fatal(fe)
		}
		return f
	}

	dir = filepath.Dir(filepath.Dir(dir))
	base := filepath.Base(dir)

	if len(base) < 5 {
		log.Fatal("Invaild dir name")
	}
	f, err := strconv.ParseFloat(base[5:], 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func GetPushData(p string) []interface{} {
	var rt []interface{}
	rt = append(rt, GetSigma(p))
	for _, v := range GetAggregateDataAll(p) {
		rt = append(rt, v.Failure)
	}
	return rt
}

// 累積和
func CumulativeSum(ads *[]AggregateData) []AggregateData {
	lines := 0
	sum := 0
	var rt []AggregateData
	for _, v := range *ads {
		lines += v.Lines
		sum += v.Failure
		rt = append(rt, AggregateData{
			Failure: sum,
			Lines:   lines,
		})
	}
	return rt
}

func PrintAggregateData(ads *[]AggregateData, failureOnly bool) {
	if failureOnly {
		for _, v := range *ads {
			fmt.Println(v.Failure)
		}
		return
	}

	for i, v := range *ads {
		fmt.Println(i, v.Lines, v.Failure)
	}
}

// s is b ?
func (s AggregateData) Compare(t AggregateData) bool {
	return s.Failure == t.Failure &&
		s.Lines == t.Lines &&
		s.FileName == t.FileName
}

// gorutine
func aggWorker(path string) (<-chan AggregateData, <-chan bool) {
	fp, err := ioutil.ReadDir(path)
	var wg sync.WaitGroup
	if err != nil {
		log.Fatal(err)
	}
	rec := make(chan AggregateData, len(fp))
	fin := make(chan bool)
	limit := make(chan struct{}, 10)
	go func() {
		for _, v := range fp {
			file := filepath.Join(path, v.Name())
			ext := filepath.Ext(file)
			if ext != ".csv" {
				continue
			}
			wg.Add(1)
			go func(p string) {
				limit <- struct{}{}
				defer wg.Done()

				res, err := NewAggregateData(p)
				log.Println(res)
				if err != nil {
					log.Fatal(err)
				}

				rec <- res
				<-limit
			}(file)
		}
		wg.Wait()
		fin <- false
	}()

	return rec, fin
}

var countFirstFilter, countSecondFilter, countThirdFilter float64

func init() {
	rootCmd.AddCommand(countCmd)

	countCmd.Flags().BoolVarP(&RangeSEEDCount, "RangeSEED", "R", false, "RangeSEEDシミュレーションの結果を数え上げます")
	countCmd.Flags().BoolP("Cumulative", "C", false, "累積和を出力します")
	countCmd.Flags().Float64Var(&countFirstFilter, "firstF", 0.4, "1カラム目のフィルターです")
	countCmd.Flags().Float64Var(&countSecondFilter, "secondF", 0.0, "2カラム目のフィルターです")
	countCmd.Flags().Float64Var(&countThirdFilter, "thirdF", 0.4, "3カラム目のフィルターです")
	countCmd.Flags().BoolP("failure-only", "F", false, "不良数とシグマだけを表示します")
}
