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

	for _, monte := range t.Simulation.Monte {
		wg.Add(1)
		cnt++

		go func(cnt int) {
			// ファイルのコピー
			//spi
			spi := getSPIScript(t.Simulation, monte)
			if err := ioutil.WriteFile(filepath.Join(t.Simulation.SimDir, "input.spi"), spi, 0644); err != nil {
				log.Fatal(err)
			}
			//ace
			ace := getACEScript(t.Simulation.Signal, t.Simulation.Range)
			if err := ioutil.WriteFile(filepath.Join(t.Simulation.DstDir, "extract.ace"), ace, 0644); err != nil {
				log.Fatal(err)
			}

			dst := filepath.Join(t.Simulation.DstDir, monte)
			tryMkdir(dst)

			command := fmt.Sprintf("cd %s &&\nhspice -hpp -mt 4 -i %s -o ./hspice &> ./hspice.log &&\nwv -k -ace_no_gui ./extract.ace &> ./wv.log &&\ncat store.csv | sed '/^#/d;1,1d' | awk -F, '{print $2}' | xargs -n3 > ../%s.csv\n", dst, filepath.Join(t.Simulation.SimDir, "input.spi"), monte)

			fmt.Println(command)

			//c := exec.Command("bash", "-c", command)

			//err := c.Run()
			//flag = flag || (err != nil)
			log.Print("Finished (", cnt, "/", size, ")")
			wg.Done()

		}(cnt)
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

func getSPIScript(s Simulation, monte string) []byte {
	return []byte(fmt.Sprintf(`.option search='%s'
.option MCBRIEF=2
.param vtn=AGAUSS(%.4f,%.4f,%.4f) vtp=AGAUSS(%.4f,%.4f,%.4f)
.option PARHIER = LOCAL
.include '%s'
.option ARTIST=2 PSF=2
.temp 25
.include '%s'
*Custom Designer (TM) Version J-2014.12-SP2-2

.GLOBAL gnd! vdd!
m30 m8d m7d vdd! vdd! PCH w=90n l=0.045u ad='(90n*0.14u)' as='(90n*0.14u)' pd='(2*(90n+0.14u))'
+  ps='(2*(90n+0.14u))'
m9 m7d m8d vdd! vdd! PCH w=90n l=0.045u ad='(90n*0.14u)' as='(90n*0.14u)' pd='(2*(90n+0.14u))'
+  ps='(2*(90n+0.14u))'
m31 blb v3 vdd! vdd! PCH1 w=90n l=0.045u ad='(90n*0.14u)' as='(90n*0.14u)' pd='(2*(90n+0.14u))'
+  ps='(2*(90n+0.14u))'
m32 bl v3 vdd! vdd! PCH1 w=90n l=0.045u ad='(90n*0.14u)' as='(90n*0.14u)' pd='(2*(90n+0.14u))'
+  ps='(2*(90n+0.14u))'
m27 m8d v2 bl gnd! NCH w=60n l=0.045u ad='(60n*0.14u)' as='(60n*0.14u)' pd='(2*(60n+0.14u))'
+  ps='(2*(60n+0.14u))'
m26 m8d m7d gnd! gnd! NCH w=60n l=0.045u ad='(60n*0.14u)' as='(60n*0.14u)' pd='(2*(60n+0.14u))'
+  ps='(2*(60n+0.14u))'
m24 bl v1 gnd! gnd! NCH1 w=300n l=0.045u ad='(300n*0.14u)' as='(300n*0.14u)' pd='(2*(300n+0.14u))'
+  ps='(2*(300n+0.14u))'
m14 blb gnd! gnd! gnd! NCH1 w=300n l=0.045u ad='(300n*0.14u)' as='(300n*0.14u)'
+ pd='(2*(300n+0.14u))' ps='(2*(300n+0.14u))'
m13 m7d v2 blb gnd! NCH w=60n l=0.045u ad='(60n*0.14u)' as='(60n*0.14u)' pd='(2*(60n+0.14u))'
+  ps='(2*(60n+0.14u))'
m25 m7d m8d gnd! gnd! NCH w=60n l=0.045u ad='(60n*0.14u)' as='(60n*0.14u)' pd='(2*(60n+0.14u))'
+  ps='(2*(60n+0.14u))'
v18 vdd! v1 dc=0 pulse ( 0.8 0 4.75n 0.5n 0.5n 9.5n 20n )
v35 vdd! v2 dc=0 pulse ( 0.8 0 4.75n 0.5n 0.5n 9.5n 20n )
v36 vdd! v3 dc=0 pulse ( 0.8 0 4.75n 0.5n 0.5n 9.5n 20n )
.tran 10p 20n start=0 uic sweep monte=%s firstrun=1
.option opfile=1 split_dp=2
.end`, s.LibDir, s.Vtn.Voltage, s.Vtn.Sigma, s.Vtn.Deviation, s.Vtp.Voltage, s.Vtp.Sigma, s.Vtp.Deviation,
		s.AddFile, s.ModelFile, monte))
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

	runCmd.PersistentFlags().Int32P("number", "n", 1, "実行するシミュレーションセットの個数です")
	runCmd.PersistentFlags().StringP("file", "f", "", "タスクファイルを指定します。一つしかできないです")
	//runCmd.PersistentFlags().Bool("fzf",false,"fzfを使ってファイルを選択します")

}
