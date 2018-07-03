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

	},
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
func copyTemplates(montes []string, dstdir string) error {
	srcdir := path.Join(os.Getenv("GOPATH"), "github.com/xztaityozx/UHA/templates/")

	mapSrc, err := os.Open(path.Join(srcdir, "resultsMap.xml"))
	defer mapSrc.Close()
	if err != nil {
		return err
	}

	for _, v := range montes {
		dist, err := os.Open(path.Join(dstdir, v, "resultsMap.xml"))
		defer dist.Close()
		if err != nil {
			return err
		}

		//Mapをコピー
		io.Copy(dist, mapSrc)

		resSrc, err := os.Open(path.Join(srcdir, v))
		defer resSrc.Close()
		if err != nil {
			return err
		}
		resDst, err := os.Open(path.Join(dstdir, v, "results.xml"))
		defer resDst.Close()
		if err != nil {
			return err
		}

		//results.xmlをコピー
		io.Copy(resDst, resSrc)

	}

	return nil
}

//func makeACEScript(monte string, signame string) string {

//}

func write(monte string, dir string, data []byte) error {
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

	fp := path.Join(p, "input.spi")

	return ioutil.WriteFile(fp, data, 0755)
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.PersistentFlags().StringP("out", "o", path.Join("/home", os.Getenv("USER"), "WorkSpace/result/"), "出力先のディレクトリです")
	makeCmd.PersistentFlags().StringP("vtn", "n", "AGAUSS(0.6,0,1.0)", "vtnの値です")
	makeCmd.PersistentFlags().StringP("vtp", "p", "AGAUSS(0.6,0,1.0)", "vtpの値です")
	makeCmd.PersistentFlags().StringSliceP("monte", "m", DEF_MOTES, "モンテカルロの回数です")
	makeCmd.PersistentFlags().StringP("signame", "s", "N2", "プロットしたい信号線の名前です")
}
