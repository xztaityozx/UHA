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
	"io/ioutil"
	"os"
	"path/filepath"

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
	},
}

type AggregateData struct {
	Lines   int
	Failure int
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
		Failure: s.State.(int),
		Lines:   s.NR,
	}, nil
}

func GetAggregateDataAll(dir string) []AggregateData {
	f, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fata(err)
	}

	for _, v := range f {
		path := filepath.Join(dir, v.Name())
		ext := filepath.Ext(path)
		if ext != "ext" {
			continue
		}

	}
}

var countFirstFilter, countSecondFilter, countThirdFilter float64

func init() {
	rootCmd.AddCommand(countCmd)

	countCmd.Flags().BoolVarP(&RangeSEEDCount, "RangeSEED", "R", false, "RangeSEEDシミュレーションの結果を数え上げます")
	countCmd.Flags().BoolP("Cumulative", "C", false, "累積和を出力します")
	countCmd.Flags().BoolP("print-sigma", "S", false, "Sigmaの値も出力します")
	countCmd.Flags().Float64Var(&countFirstFilter, "firstF", 0.4, "1カラム目のフィルターです")
	countCmd.Flags().Float64Var(&countSecondFilter, "secondF", 0.0, "2カラム目のフィルターです")
	countCmd.Flags().Float64Var(&countThirdFilter, "thirdF", 0.4, "3カラム目のフィルターです")
}
