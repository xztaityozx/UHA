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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func getFromDst(src string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.Name()[0:5] != "Sigma" || !f.IsDir() {
			continue
		}

		s := filepath.Join(src, f.Name())
		if len(config.Repository) == 0 {
			return errors.New("Could not find Repository")
		}

		if config.Repository[0].Type != Git && config.Repository[0].Type != Dir {
			return errors.New("いまはGitかディレクトリしか操作できません")
		}

		t := filepath.Join(config.Repository[0].Path, f.Name())
		if err := os.Rename(s, t); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().BoolP("push", "P", false, "SpreadSheetにもデータを書き込みます")
	getCmd.PersistentFlags().StringP("from", "f", config.Simulation.DstDir, fmt.Sprint("Sigmax.xxがあるフォルダです(default ", config.Simulation.DstDir, ")"))

}
