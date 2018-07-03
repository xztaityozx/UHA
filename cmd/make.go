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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		vtn, err := cmd.PersistentFlags().GetString("vtn")
		if err != nil {
			log.Fatal(err)
		}
		vtp, err := cmd.PersistentFlags().GetString("vtp")
		if err != nil {
			log.Fatal(err)
		}

		ml, err := cmd.PersistentFlags().GetStringSlice("monte")
		if err != nil {
			log.Fatal(err)
		}

		dir, err := cmd.PersistentFlags().GetString("out")
		if err != nil {
			log.Fatal(err)
		}

		r, err := cmd.PersistentFlags().GetStringArray("range")
		if err != nil || len(r) != 3 {
			log.Fatal("Fatal range")
		}

		signame, err := cmd.PersistentFlags().GetString("signame")
		if err != nil {
			log.Fatal(err)
		}

		if err := FireTask(ml, vtp, vtn, Range{Start: r[0], Stop: r[1], Step: r[2]}, dir, signame); err != nil {
			log.Fatal(err)
		}

	},
}

func FireTask(ms []string, vtp string, vtn string, r Range, dst string, signame string) error {
	log.Printf("Starting Simulation Set\n\tVtp = %s\n\tVtn = %s\n\tRangeStart = %s\n\tRangeStop = %s\n\tRangeStep = %s\n", vtp, vtn, r.Start, r.Stop, r.Step)

	cnt := 0
	index := 0

	wg := new(sync.WaitGroup)

	s := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	log.Print("start simulation")
	s.Suffix = fmt.Sprintf(" Dispatch Tasks")
	s.Writer = os.Stderr
	s.Start()
	defer s.Stop()
	for _, v := range ms {
		wg.Add(1)
		go func(cnt int) {

			// write SPI/ACE script
			if err := write(v, dst, []byte(makeSPI(vtn, vtp, v)), []byte(makeACEScript(signame, r.Start, r.Stop, r.Step))); err != nil {
				log.Fatal(err)
			}

			// copyTemplate to dist dir
			if err := copyTemplate(v, dst); err != nil {
				log.Fatal(err)
			}

			// start simulation
			fmt.Printf("cd %s && hspice -hpp -mt 4 ./input.spi > ./hspice.log && wv -ace_no_gui ./extract.ace -k > ./wv.log && ", path.Join(dst, v))
			fmt.Printf("cat store.csv | sed '1,1d;/^#/d'|awk -F, '{print $2}'|xargs -n3 > %s.csv && ", v)
			fmt.Printf("cat %s.csv | awk '$1>=0.4&&$3>=0.4{print}'| wc -l\n", v)

			time.Sleep(time.Duration(cnt+1) * time.Second)

			index++
			log.Print(" Dispatch Tasks (", index, "/", len(ms), ")")
			wg.Done()
		}(cnt)
		cnt++
	}
	s.FinalMSG = "All Tasks Ended"
	wg.Wait()
	return nil
}

type Range struct {
	Start string
	Stop  string
	Step  string
}

func makeSPI(vtn string, vtp string, monte string) string {
	return fmt.Sprintf(`*  Generated for: HSPICE
*  Design library name: takemura
*  Design cell name: sram
*  Design view name: sram
.option search='/home/takemura/Workspace/takemura'

.option MCBRIEF=2
.param vtn=%s vtp=%s
.option PARHIER = LOCAL
.include '/home/takemura/Workspace/takemura/addfile.txt'
.option ARTIST=2 PSF=2
.temp 25
.include 'modified02_45nm_bulk_BSIM4_v1.0_HSPICE_pm.txt'
*Custom Designer (TM) Version J-2014.12-SP2-2
*Mon Jul  2 16:21:08 2018

.GLOBAL gnd! vdd!
********************************************************************************
* Library          : takemura
* Cell             : sram
* View             : sram
* View Search List : hspice hspiceD schematic spice veriloga
* View Stop List   : hspice hspiceD
********************************************************************************
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




.end`, vtn, vtp, monte)
}

//テンプレートをコピーするやつ
func copyTemplate(monte string, dstdir string) error {
	srcdir := path.Join(os.Getenv("GOPATH"), "src/github.com/xztaityozx/UHA/templates/")

	mapSrc, err := os.Open(path.Join(srcdir, "resultsMap.xml"))
	defer mapSrc.Close()
	if err != nil {
		return err
	}

	mapDst, err := os.OpenFile(path.Join(dstdir, monte, "resultsMap.xml"), os.O_CREATE|os.O_WRONLY, 0644)
	defer mapDst.Close()
	if err != nil {
		return err
	}

	if _, err := io.Copy(mapDst, mapSrc); err != nil {
		return err
	}

	resSrc, err := os.Open(path.Join(srcdir, monte))
	defer resSrc.Close()
	if err != nil {
		return err
	}

	resDst, err := os.OpenFile(path.Join(dstdir, monte, "results.xml"), os.O_CREATE|os.O_WRONLY, 0644)
	defer resDst.Close()
	if err != nil {
		return err
	}

	if _, err := io.Copy(resDst, resSrc); err != nil {
		return err
	}
	return nil
}

func makeACEScript(signame string, start string, stop string, step string) string {
	return fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
set www [ sx_find_wave_in_file $xml %s ]
sx_export_csv on
sx_export_range %s %s %s
sx_export_data  "store.csv" $www
	`, signame, start, stop, step)
}

func write(monte string, dir string, spi []byte, ace []byte) error {
	if _, err := os.Stat(dir); err != nil {
		log.Print("make ", dir)
		if e := os.Mkdir(dir, 0755); e != nil {
			return e
		}
	}

	p := path.Join(dir, monte)
	if _, err := os.Stat(p); err != nil {
		log.Print("make ", p)
		if e := os.Mkdir(p, 0755); e != nil {
			return e
		}
	}

	fspi := path.Join(p, "input.spi")
	if err := ioutil.WriteFile(fspi, spi, 0644); err != nil {
		return err
	}

	face := path.Join(p, "extract.ace")
	return ioutil.WriteFile(face, ace, 0644)
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.PersistentFlags().StringP("out", "o", path.Join("/home", os.Getenv("USER"), "WorkSpace/result/"), "出力先のディレクトリです")
	makeCmd.PersistentFlags().StringP("vtn", "n", "AGAUSS(0.6,0,1.0)", "vtnの値です")
	makeCmd.PersistentFlags().StringP("vtp", "p", "AGAUSS(0.6,0,1.0)", "vtpの値です")
	makeCmd.PersistentFlags().StringSliceP("monte", "m", DEF_MOTES, "モンテカルロの回数です")
	makeCmd.PersistentFlags().StringP("signame", "s", "N2", "プロットしたい信号線の名前です")
	makeCmd.PersistentFlags().StringArrayP("range", "r", []string{"2.5ns", "17.5ns", "7.5ns"}, "時間を指定します")
}
