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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	RESERVE string = "reserve"
	DONE    string = "done"
	FAILED  string = "failed"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "シミュレーションを実行します",
	Long:  `シミュレーションセットを実行します`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func runTask(t Task, f string) error {

}

func readTask() (Task, string) {
	p := config.TaskDir

	// リスト取得
	files, err := ioutil.ReadDir(filepath.Join(p, RESERVE))
	if err != nil {
		log.Fatal(err)
	}

	f := filepath.Join(p, RESERVE, files[0].Name())

	//実行と移動
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	var task Task
	err := json.Unmarshal(b, &task)

	if err != nil {
		moveTo(files[0].Name(), FAILED)
		log.Fatal(err)
	}

	return task
}

func tryMkdir(p string) error {
	if _, err := os.Stat(p); err != nil {
		if e := os.MkdirAll(p, 0755); e != nil {
			return e
		}
		log.Print("Mkdir : ", p)
	}
}

func moveTo(f string, dir string) {
	src := filepath.Join(ReserveDir, f)
	r, err := os.OpenFile(src, os.O_CREATE|os.O_RDONLY, 0644)
	defer r.Close()
	if err != nil {
		log.Fatal(err)
	}

	dst := filepath.Join(dir, f)
	w, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	defer w.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Move to ", dst)
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().Int32P("number", "n", 1, "実行するシミュレーションセットの個数です")
	runCmd.PersistentFlags().StringP("file", "f", "", "タスクファイルを指定します。一つしかできないです")
	//runCmd.PersistentFlags().Bool("fzf",false,"fzfを使ってファイルを選択します")

}
