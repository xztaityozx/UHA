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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/briandowns/spinner"
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
		t, f := readTask()
		if err := runTask(t); err != nil {
			moveTo(f, FailedDir)
			log.Fatal(err)
		} else {
			moveTo(f, DoneDir)
		}
	},
}

func runTask(t Task) error {
	s := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " Running... "
	s.Start()

	var wg sync.WaitGroup

	flag := false
	cnt := 0
	size := len(t.Simulation.Monte)

	outDir := filepath.Join(t.Simulation.DstDir, fmt.Sprintf("Sigma%s", t.Simulation.Vtn.Sigma))
	if err := tryMkdir(outDir); err != nil {
		return err
	}

	for _, monte := range t.Simulation.Monte {
		wg.Add(1)
		cnt++

		go func(monte string, cnt int) {
			defer wg.Done()
			dst := filepath.Join(t.Simulation.DstDir, monte)
			tryMkdir(dst)
			// ファイルのコピー
			//spi
			spi, re := getSPIScript(t.Simulation, monte)
			if re != nil {
				log.Println(re)
				flag = true
				return
			}
			if err := ioutil.WriteFile(filepath.Join(t.Simulation.SimDir, fmt.Sprintf("%sinput.spi", monte)), spi, 0644); err != nil {
				log.Println(err)
				flag = true
				return
			}
			//ace
			ace := getACEScript(t.Simulation.Signal, t.Simulation.Range)
			if err := ioutil.WriteFile(filepath.Join(t.Simulation.DstDir, "extract.ace"), ace, 0644); err != nil {
				log.Println(err)
				flag = true
				return
			}

			//resultMap
			rmap, rmaperr := ioutil.ReadFile(filepath.Join(SelfPath, "templates", "resultsMap.xml"))
			if rmaperr != nil {
				log.Println(rmaperr)
				flag = true
				return
			}
			if err := ioutil.WriteFile(filepath.Join(t.Simulation.DstDir, monte, "resultsMap.xml"), rmap, 0644); err != nil {
				log.Println(err)
				flag = true
				return
			}
			//result
			res, reserr := ioutil.ReadFile(filepath.Join(SelfPath, "templates", monte))
			if reserr != nil {
				log.Println(reserr)
				flag = true
				return
			}
			if err := ioutil.WriteFile(filepath.Join(t.Simulation.DstDir, monte, "results.xml"), res, 0644); err != nil {
				log.Println(err)
				flag = true
				return
			}

			command := fmt.Sprintf("cd %s &&\nhspice -hpp -mt 4 -i %s -o ./hspice &> ./hspice.log &&\nwv -k -ace_no_gui ../extract.ace &> ./wv.log &&\ncat store.csv | sed '/^#/d;1,1d' | awk -F, '{print $2}' | xargs -n3 > ../Sigma%.4f/%s.csv\n", dst, filepath.Join(t.Simulation.SimDir, fmt.Sprintf("%sinput.spi", monte)), t.Simulation.Vtn.Sigma, monte)

			//fmt.Println(command)

			c := exec.Command("bash", "-c", command)

			err := c.Run()
			flag = flag || (err != nil)
			log.Print("Finished (", cnt, "/", size, ")")

		}(monte, cnt)
	}
	wg.Wait()

	s.Stop()
	s.FinalMSG = "simulation has finished"

	if flag {
		return errors.New("simulation set has failed")
	}

	return nil
}

func getACEScript(s string, r Range) []byte {
	return []byte(fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
	set www [ sx_find_wave_in_file $xml %s ]
	sx_export_csv on
	sx_export_range %s %s %s
	sx_export_data  "store.csv" $www
	`, s, r.Start, r.Stop, r.Step))
}

func getSPIScript(s Simulation, monte string) ([]byte, error) {
	p := filepath.Join(ConfigDir, "spitemplate.spi")
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return []byte{}, err
	}
	tmplt := string(b)
	return []byte(fmt.Sprintf(tmplt, s.Vtn.Voltage, s.Vtn.Sigma, s.Vtn.Deviation,
		s.Vtp.Voltage, s.Vtp.Sigma, s.Vtp.Deviation, monte)), nil
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
	err = json.Unmarshal(b, &task)

	if err != nil {
		moveTo(files[0].Name(), FAILED)
		log.Fatal(err)
	}

	return task, files[0].Name()
}

func tryMkdir(p string) error {
	if _, err := os.Stat(p); err != nil {
		if e := os.MkdirAll(p, 0755); e != nil {
			return e
		}
		log.Print("Mkdir : ", p)
	}
	return nil
}

func moveTo(f string, dir string) {
	src := filepath.Join(ReserveDir, f)

	dst := filepath.Join(dir, f)

	if err := os.Rename(src, dst); err != nil {
		log.Fatal(err)
	}

	log.Print("Move to ", dst)
}

func init() {
	rootCmd.AddCommand(runCmd)

	//runCmd.PersistentFlags().Int32P("number", "n", 1, "実行するシミュレーションセットの個数です")
	//runCmd.PersistentFlags().StringP("file", "f", "", "タスクファイルを指定します。一つしかできないです")
	//runCmd.PersistentFlags().Bool("fzf",false,"fzfを使ってファイルを選択します")
}
