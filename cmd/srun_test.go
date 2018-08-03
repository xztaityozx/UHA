package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestAllSRun(t *testing.T) {
	var home string = os.Getenv("HOME")
	//var taskDir string = filepath.Join(home, ".config", "UHA", "task")
	ConfigDir = filepath.Join(home, ".config", "UHA")
	SelfPath = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "xztaityozx", "UHA")
	var rst RangeSEEDTask = RangeSEEDTask{
		BaseDir: filepath.Join(home, "WorkSpace", "Base"),
		Sim:     filepath.Join(home, "WorkSpace", "Sim"),
		Monte:   "500",
		SEED:    1,
		Vtn: Node{
			Voltage:   0.6,
			Sigma:     0.1,
			Deviation: 1.0,
		},
		Vtp: Node{
			Voltage:   -0.6,
			Sigma:     0.1,
			Deviation: 1.0,
		},
		Sigma: 0.1,
	}

	tryMkdir(rst.Sim)

	t.Run("001 tryMkRangeSEEDDstDir", func(t *testing.T) {
		if err := tryMkRangeSEEDDstDir(&rst); err != nil {
			t.Fatal(err)
		}

		p := filepath.Join(rst.BaseDir, fmt.Sprintf("RangeSEED_Vtn%.4fVtp%.4f_Sigma%.4f_Monte%s/SEED%03d", rst.Vtn.Sigma, rst.Vtp.Sigma, rst.Sigma, rst.Monte, rst.SEED))
		if _, err := os.Stat(p); err != nil {
			t.Fatal(err)
		}

		if rst.Dst != p {
			t.Fatal("Unexpected rst.Dst")
		}
	})

	t.Run("002 writeRangeSEEAddfile", func(t *testing.T) {
		if err := writeRangeSEEAddfile(&rst); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(filepath.Join(rst.BaseDir, "Addfiles")); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(filepath.Join(rst.BaseDir, "Addfiles", "addfile001.txt")); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("003 writeRangeSEEDSPI", func(t *testing.T) {
		if err := writeRangeSEEDSPI(&rst); err != nil {
			t.Fatal(err)
		}

		f := filepath.Join(rst.Sim, "Vtn0.1000Vtp0.1000Monte500_Sigma0.1000_SEED001.spi")
		if _, err := os.Stat(f); err != nil {
			t.Fatal(err)
		}

		if f != rst.SPI {
			t.Fatal("Unexpected rst.SPI")
		}
	})

	t.Run("004 makeRangeSEEDCommand", func(t *testing.T) {
		expect := fmt.Sprintf("cd %s && hspice -mt 2 -i %s -o ./hspice &> ./hspice.log &&", rst.Dst, rst.SPI)
		expect += fmt.Sprintf("wv -k -ace_no_gui ../../extract.ace &> ./wv.log &&")
		expect += fmt.Sprintf("cat store.csv | sed '/^#/d;1,1d' | awk -F, '{print $2}' | xargs -n3 > ../Result/SEED%03d.csv ", rst.SEED)

		log.Println(expect)

		actual := makeRangeSEEDCommand(rst)
		if actual != expect {
			t.Fatal("Unexpected result : ", actual)
		}
	})

	t.Run("005 copyRangeSEEDXmls", func(t *testing.T) {
		if err := copyRangeSEEDXmls(rst); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("006 rmeove", func(t *testing.T) {
		os.RemoveAll(rst.Dst)
		os.RemoveAll(rst.Sim)
		os.RemoveAll(rst.BaseDir)
	})
}
