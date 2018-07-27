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

	"github.com/spf13/cobra"
	"github.com/wayneashleyberry/terminal-dimensions"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "タスクをリストアップします",
	Long: `UHAの管理下にあるタスクをリストアップします
Usage:
	UHA list [--one,-1|--long,-l]

読み込んだコンフィグで設定されている"TaskDir"以下にあるタスクをリストアップします

`,
	Run: func(cmd *cobra.Command, args []string) {
		var one, long bool
		var err error

		one, err = cmd.PersistentFlags().GetBool("one")
		if err != nil {
			log.Fatal(err)
		}

		long, err = cmd.PersistentFlags().GetBool("long")
		if err != nil {
			log.Fatal(err)
		}

		var res []string

		var items []listDirectory
		items, err = listReadAllItems()
		if err != nil {
			log.Fatal(err)
		}
		if one {
			res = append(res, ListPrefix[3])
			res = append(res, singleLineList(&items[0])...)
			res = append(res, ListPrefix[4])
			res = append(res, singleLineList(&items[1])...)

		} else if long {
			res = append(res, ListPrefix[3])
			res = append(res, longList(&items[0])...)
			res = append(res, ListPrefix[4])
			res = append(res, longList(&items[1])...)

		} else {
			w, _ := terminaldimensions.Width()

			res = append(res, ListPrefix[3])
			res = append(res, multiLineList(&items[0], int(w))...)
			res = append(res, ListPrefix[4])
			res = append(res, multiLineList(&items[1], int(w))...)
		}

		for _, v := range res {
			fmt.Println(v)
		}
	},
}

type listItem struct {
	Name string
	Date string
}

type listDirectory struct {
	Reserve []listItem
	Done    []listItem
	Failed  []listItem
}

// 各タスクDirからすべて読み出す
func listReadAllItems() ([]listDirectory, error) {
	// Run
	run := listDirectory{
		Reserve: listReadDir(ReserveRunDir),
		Done:    listReadDir(DoneRunDir),
		Failed:  listReadDir(FailedRunDir),
	}

	// SRun
	srun := listDirectory{
		Reserve: listReadDir(ReserveSRunDir),
		Done:    listReadDir(DoneSRunDir),
		Failed:  listReadDir(FailedSRunDir),
	}

	return []listDirectory{run, srun}, nil
}

// p以下を読み取って、listItemの配列にして返す
func listReadDir(p string) []listItem {
	if files, err := ioutil.ReadDir(p); err != nil {
		log.Fatal(err)
		return nil
	} else {
		var rt []listItem
		for _, f := range files {
			l := listItem{
				Name: f.Name(),
				Date: f.ModTime().Format("2006/01/02/15:04"),
			}
			rt = append(rt, l)
		}
		return rt
	}
}

// 複数カラムでリストアップします。
// list : リストアップ対象
// width : ターミナルの幅
func multiLineList(list *listDirectory, width int) []string {
	var rt []string
	length := 0
	if len(list.Reserve) != 0 {
		length = len(list.Reserve[0].Name)
	} else if len(list.Done) != 0 {
		length = len(list.Done[0].Name)
	} else if len(list.Failed) != 0 {
		length = len(list.Failed[0].Name)
	} else {
		return []string{}
	}
	column := width / (length + 4)

	// Reserve
	rt = append(rt, fmt.Sprintf(ListPrefix[0], len(list.Reserve)))
	rt = append(rt, getMultiLine(list.Reserve, column)...)
	rt = append(rt, "")

	// Done
	rt = append(rt, fmt.Sprintf(ListPrefix[1], len(list.Done)))
	rt = append(rt, getMultiLine(list.Done, column)...)
	rt = append(rt, "")

	// Failed
	rt = append(rt, fmt.Sprintf(ListPrefix[2], len(list.Failed)))
	rt = append(rt, getMultiLine(list.Failed, column)...)
	rt = append(rt, "")
	return rt
}

var ListPrefix []string = []string{
	"\033[1;33m●\033[0;39m  Reserve: %d個",
	"\033[1;32m●\033[0;39m  Done: %d個",
	"\033[1;31m●\033[0;39m  Failed: %d個",
	"\033[1;32m>>\033[0;39m Task of Run:",
	"\033[1;32m>>\033[0;39m Task of SRun:"}

// 複数カラムのリストアップ用の1行を作ります
func getMultiLine(l []listItem, c int) []string {
	var rt []string
	cnt := 0
	line := ""

	for _, v := range l {
		line += fmt.Sprintf("\t%s", v.Name)
		cnt++
		if cnt == c {
			rt = append(rt, line)
			cnt = 0
			line = ""
		}
	}
	if len(line) != 0 {
		rt = append(rt, line)
	}
	return rt
}

// ファイル名と時間の組み合わせからなるLongなリストアップをします
func longList(list *listDirectory) []string {
	var rt []string
	// Reserve
	rt = append(rt, fmt.Sprintf(ListPrefix[0], len(list.Reserve)))
	rt = append(rt, getLongList(list.Reserve)...)
	rt = append(rt, "")

	// Done
	rt = append(rt, fmt.Sprintf(ListPrefix[1], len(list.Done)))
	rt = append(rt, getLongList(list.Done)...)
	rt = append(rt, "")

	// Failed
	rt = append(rt, fmt.Sprintf(ListPrefix[2], len(list.Failed)))
	rt = append(rt, getLongList(list.Failed)...)
	rt = append(rt, "")

	return rt
}

// Longの1行を作ります
func getLongList(li []listItem) []string {
	var rt []string
	tr := len("20180727153226_")
	for _, v := range li {
		if len(v.Name) > tr {
			rt = append(rt, fmt.Sprintf("%s\t%s", v.Name[tr:], v.Date))
		} else {
			rt = append(rt, fmt.Sprintf("%s\t%s", v.Name, v.Date))
		}

	}
	return rt
}

// 1カラムだけでリストアップします
func singleLineList(list *listDirectory) []string {
	var rt []string
	// Reserve
	rt = append(rt, fmt.Sprintf(ListPrefix[0], len(list.Reserve)))
	rt = append(rt, getSingleLineList(list.Reserve)...)
	rt = append(rt, "")

	// Done
	rt = append(rt, fmt.Sprintf(ListPrefix[1], len(list.Done)))
	rt = append(rt, getSingleLineList(list.Done)...)
	rt = append(rt, "")

	// Failed
	rt = append(rt, fmt.Sprintf(ListPrefix[2], len(list.Failed)))
	rt = append(rt, getSingleLineList(list.Failed)...)
	rt = append(rt, "")

	return rt
}

// 1カラムの1行を作ります
func getSingleLineList(li []listItem) []string {
	var rt []string
	for _, v := range li {
		rt = append(rt, v.Name)
	}
	return rt
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolP("long", "l", false, "詳細を表示します")
	listCmd.PersistentFlags().BoolP("one", "1", false, "強制的に1行に1つ出力します")
}
