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

	"github.com/spf13/cobra"
)

// srunCmd represents the srun command
var srunCmd = &cobra.Command{
	Use:   "srun",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("srun called")
	},
}

// ReserveSRunDirから、NSeedTaskのJSONとして正しいやつだけ列挙
func readNSTaskFileList() []NSeedTask {
	var rt []NSeedTask
	files, err := ioutil.ReadDir(ReserveSRunDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		p := filepath.Join(ReserveSRunDir, f.Name())
		b, rerr := ioutil.ReadFile(p)
		if rerr != nil {
			log.Fatal(rerr)
		}
		var nt NSeedTask
		j, jerr := json.Unmarshal(b, &nt)
		if jerr != nil {
			log.Println(jerr)
		}

		rt = append(rt, nt)
	}

	return rt
}

func setResultDir(nt NSeedTask) error {
	for i := 1; i <= nt.Count; i++ {
		p := filepath.Join()
	}
}

func makeSRunCommand(nt NSeedTask) []string {
	var rt []string
}

func setAddfileTo(cnt int, p string) error {
	for i := 1; i <= cnt; i++ {
		b, err := []byte(makeSAddfile(i))
		if err != nil {
			return err
		}

		f := filepath.Join(p, fmt.Sprintf("addfile%d.txt", i))
		if err := ioutil.WriteFile(f, b, 0644); err != nil {
			return err
		}
		log.Printf("Write Addfile%d To : %s\n", i, f)
	}
	return nil
}

// ConfigDir以下にあるaddfile.txtをテンプレートに、SEEDを変更したaddfileの文字列を作る
func makeSAddfile(seed int) (string, error) {
	p := filepath.Join(ConfigDir, "addfile.txt")
	tmp, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(string(tmp), seed), nil
}

func init() {
	rootCmd.AddCommand(srunCmd)
	srunCmd.PersistentFlags().Bool("all", false, "すべて実行します")
}
