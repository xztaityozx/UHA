package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestAllRun(t *testing.T) {
	home := os.Getenv("HOME")
	ConfigDir = filepath.Join(home, ".config", "UHA")
	SelfPath = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "xztaityozx", "UHA")
	ReserveRunDir = filepath.Join(ConfigDir, "task", "run", RESERVE)
	tryMkdirSuppress(ReserveRunDir)

	rt := RunTask{
		Simulation: Simulation{
			DstDir: filepath.Join(home, "WorkSpace", "Dst"),
			SimDir: filepath.Join(home, "WorkSpace", "Sim"),
			Monte:  []string{"100", "200", "500"},
			Range: Range{
				Start: "0ns",
				Stop:  "10ns",
				Step:  "0.5ns",
			},
			Signal: "N2",
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
			SEED: 1,
		},
		SEED: 1,
	}

	t.Run("000_Prepare", func(t *testing.T) {
		addfile := filepath.Join(ConfigDir, "addfile.txt")
		if err := ioutil.WriteFile(addfile, []byte("%d"), 0644); err != nil {
			t.Fatal(err)
		}

		spi := filepath.Join(ConfigDir, "spitemplate.spi")
		if err := ioutil.WriteFile(spi, []byte("%.4f\n%.4f\n%.4f\n%.4f\n%.4f\n%.4f\n%s\n%s"), 0644); err != nil {
			t.Fatal(err)
		}

		if err := tryMkdirSuppress(rt.Simulation.SimDir); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("001_tryMkRunDstDir", func(t *testing.T) {
		if err := tryMkRunDstDir(&rt); err != nil {
			t.Fatal(err)
		}

		baseExpect := filepath.Join(home, "WorkSpace", "Dst", "VtpVolt-0.6000_VtnVolt0.6000")
		if baseExpect != rt.Base {
			t.Fatal("Unexpected Result Base : ", rt.Base, "\nexpect : ", baseExpect)
		}

		dstExpect := []string{
			filepath.Join(baseExpect, "SEED001/Monte100"),
			filepath.Join(baseExpect, "SEED001/Monte200"),
			filepath.Join(baseExpect, "SEED001/Monte500"),
		}
		for i, v := range dstExpect {
			if v != rt.Dst[i] {
				t.Fatal("Unexpected Result Dst[", i, "] : ", rt.Dst[i], "\nexpect : ", v)
			}
		}

		rfExpect := []string{
			filepath.Join(baseExpect, "SEED001/Result/100.csv"),
			filepath.Join(baseExpect, "SEED001/Result/200.csv"),
			filepath.Join(baseExpect, "SEED001/Result/500.csv"),
		}
		for i, v := range rfExpect {
			if v != rt.ResultFile[i] {
				t.Fatal("Unexpected Result Dst[", i, "] : ", rt.ResultFile[i], "\nexpect : ", v)
			}
		}
	})

	t.Run("002_tryMkRunAddfile", func(t *testing.T) {
		if err := tryMkRunAddfile(&rt); err != nil {
			t.Fatal(err)
		}

		expect := filepath.Join(rt.Base, "Addfiles", "RunAddfileSEED001.txt")
		actual := rt.Addfile

		if expect != actual {
			t.Fatal("Unexpected Result Addfile : ", actual, "\nexpect : ", expect)
		}

		addfile, err := ioutil.ReadFile(expect)
		if err != nil {
			t.Fatal(err)
		}

		addfileActual := string(addfile)
		addfileExpect := "1"

		if addfileActual != addfileExpect {
			t.Fatal("Unexpected Result Addfile : ", addfileActual, "\nexpect : ", addfileExpect)
		}

	})

	t.Run("003_tryMkRunSPI", func(t *testing.T) {
		if err := tryMkRunSPI(&rt); err != nil {
			t.Fatal(err)
		}

		expect := []string{
			filepath.Join(rt.Simulation.SimDir, "VtpVolt-0.6000_VtnVolt0.6000_SEED001_Monte100.spi"),
			filepath.Join(rt.Simulation.SimDir, "VtpVolt-0.6000_VtnVolt0.6000_SEED001_Monte200.spi"),
			filepath.Join(rt.Simulation.SimDir, "VtpVolt-0.6000_VtnVolt0.6000_SEED001_Monte500.spi"),
		}
		actual := rt.SPI

		for i, v := range actual {
			if v != expect[i] {
				t.Fatal("Unexpected Result SPI[", i, "] :", v, "\nexpect : ", expect[i])
			}
		}

		spiExpect := []string{
			fmt.Sprintf("0.6000\n0.1000\n1.0000\n-0.6000\n0.1000\n1.0000\n%s\n100", rt.Addfile),
			fmt.Sprintf("0.6000\n0.1000\n1.0000\n-0.6000\n0.1000\n1.0000\n%s\n200", rt.Addfile),
			fmt.Sprintf("0.6000\n0.1000\n1.0000\n-0.6000\n0.1000\n1.0000\n%s\n500", rt.Addfile),
		}

		for i, v := range spiExpect {
			get, err := ioutil.ReadFile(expect[i])
			if err != nil {
				t.Fatal(err)
			}

			if string(get) != v {
				t.Fatal("Unexpected Result spiExpect[", i, "] : ", string(get), "\nexpect", v)
			}
		}

	})

	t.Run("004_tryCopyRunXmls", func(t *testing.T) {
		if err := tryCopyRunXmls(rt); err != nil {
			t.Fatal(err)
		}

		expects := []string{
			filepath.Join(rt.Dst[0], "resultsMap.xml"),
			filepath.Join(rt.Dst[0], "results.xml"),
			filepath.Join(rt.Dst[1], "resultsMap.xml"),
			filepath.Join(rt.Dst[1], "results.xml"),
			filepath.Join(rt.Dst[2], "resultsMap.xml"),
			filepath.Join(rt.Dst[2], "results.xml"),
		}
		for _, v := range expects {
			if _, err := os.Stat(v); err != nil {
				t.Fatal("Unexpected Result Xmls : Not found :", v)
			}
		}
	})
	t.Run("005 ExtractFromStoreCSV", func(t *testing.T) {
		text := []byte(`# head1
# head2
TIME ,Signal
#sweep 1
 1.0000EXX , 2.0000EXX
 3.0000EXX , 4.0000EXX
 5.0000EXX , 6.0000EXX
#sweep 2
 7.0000EXX , 8.0000EXX
 9.0000EXX , 1.0000EXX
 1.1000EXX , 1.2000EXX
`)

		for _, v := range rt.Dst {
			if err := ioutil.WriteFile(filepath.Join(v, "store.csv"), text, 0644); err != nil {
				t.Fatal(err)
			}
		}

		expect := "2.0000EXX 4.0000EXX 6.0000EXX\n8.0000EXX 1.0000EXX 1.2000EXX"

		for i, csv := range rt.ResultFile {
			if err := ExtractFromStoreCSV(rt.Dst[i], csv); err != nil {
				t.Fatal(err)
			}
			if b, err := ioutil.ReadFile(csv); err != nil {
				t.Fatal(err)
			} else {
				if expect != string(b) {
					t.Fatal("Unexpected Result CSV : ", string(b), "\nexpect : ", expect)
				}
			}
		}
	})

	t.Run("006_tryMkRunACE", func(t *testing.T) {
		if err := tryMkRunACE(&rt); err != nil {
			t.Fatal(err)
		}

		pathExpect := filepath.Join(rt.Base, "extract.ace")
		if pathExpect != rt.ACE {
			t.Fatal("Unexpected Result Ace Path : ", rt.ACE, "\nexpect : ", pathExpect)
		}

		expect := getACEScript(rt.Simulation.Signal, rt.Simulation.Range)
		actual, err := ioutil.ReadFile(pathExpect)
		if err != nil {
			t.Fatal(err)
		}

		if string(actual) != string(expect) {
			t.Fatal("Unexpected Result ACE")
		}
	})

	t.Run("007_mkRunCommand", func(t *testing.T) {
		actual, err := mkRunCommand(rt)
		if err != nil {
			t.Fatal(err)
		}

		expects := []string{
			fmt.Sprintf("cd %s && hspice -mt 2 -i %s -o ./hspice &> ./hspice.log && wv -k -ace_no_gui", rt.Dst[0], rt.SPI[0]),
			fmt.Sprintf("cd %s && hspice -mt 2 -i %s -o ./hspice &> ./hspice.log && wv -k -ace_no_gui", rt.Dst[1], rt.SPI[1]),
			fmt.Sprintf("cd %s && hspice -mt 2 -i %s -o ./hspice &> ./hspice.log && wv -k -ace_no_gui", rt.Dst[2], rt.SPI[2]),
		}

		for i, v := range expects {
			if v != actual[i] {
				t.Fatal("Unexpected Result Command[", i, "] : ", actual[i], "\nexpect : ", v)
			}
		}

		if _, err := mkRunCommand(RunTask{}); err == nil {
			t.Fatal("Unexpected Result NullTask")
		}
	})

	t.Run("008_removeRunGarbage", func(t *testing.T) {
		if err := removeRunGarbage(rt); err != nil {
			t.Fatal(err)
		}

		sim, _ := ioutil.ReadDir(rt.Simulation.SimDir)
		if len(sim) != 0 {
			t.Fatal("Unexpected Result SPI : ", len(sim))
		}

		dir := []string{
			filepath.Join(rt.Base, "SEED001/Monte100"),
			filepath.Join(rt.Base, "SEED001/Monte200"),
			filepath.Join(rt.Base, "SEED001/Monte500"),
		}
		for _, v := range dir {
			if _, err := os.Stat(v); err == nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("009_MakeRunTask", func(t *testing.T) {
		b, _ := json.Marshal(rt.Simulation)
		f := filepath.Join(ConfigDir, "task", "run", RESERVE, "TEST.json")
		ioutil.WriteFile(f, b, 0644)

		if res, err := MakeRunTask(f); err != nil {
			log.Fatal(err)
		} else {
			if res.SEED != 1 {
				t.Fatal("Unexpected Result SEED : ", res.SEED)
			}
			if res.Simulation.DstDir != rt.Simulation.DstDir {
				t.Fatal("Unexpected Result DstDir : ", res.Simulation.DstDir)
			}
			if res.TaskFile != f {
				t.Fatal("Unexpected Result TaskFile : ", res.TaskFile, "\nexpect : ", f)
			}
		}
	})

	t.Run("010_readRunTasks", func(t *testing.T) {
		res, err := readRunTasks(1, false)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("Unexpected Result len(res) : ", len(res))
		}

		b, _ := json.Marshal(rt.Simulation)
		f := filepath.Join(ConfigDir, "task", "run", RESERVE, "TEST1.json")
		ioutil.WriteFile(f, b, 0644)

		res, err = readRunTasks(2, false)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("Unexpected Result len(res) : ", len(res))
		}

		res, err = readRunTasks(1, true)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("Unexpected Result len(res) : ", len(res))
		}

		for _, v := range res {
			os.Remove(v.TaskFile)
		}
	})

}
