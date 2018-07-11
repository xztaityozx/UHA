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
	"strings"

	"github.com/mattn/go-pipeline"
	"github.com/spakin/awk"
	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		Count()
	},
}

func Count() {
	rj := readPushData()
	rj.Data = aggregate()
	writePushData(rj)
}

func countup(p string) (int, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return -1, err
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
		return -1, err
	}

	return s.State.(int), nil

}

func aggregate() []interface{} {
	var rt []interface{}

	wd, _ := os.Getwd()
	if len(wd) < len("Sigmax.xxxx") {
		log.Fatal("カレントディレクトリの命名ルールが違います")
	}
	wl := len(wd)
	sigma := wd[wl-6 : wl-1]
	log.Print("Open Sigma : ", sigma)

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
		cnt, err := countup(filepath.Join(wd, v))
		if err != nil {
			log.Fatal(err)
		}

		rt = append(rt, cnt)
		fmt.Println(v, cnt)
	}
	return rt
}

func init() {
	rootCmd.AddCommand(countCmd)
}
