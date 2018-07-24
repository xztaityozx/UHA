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

var home string = os.Getenv("HOME")
var taskDir string = filepath.Join(home, ".config", "UHA", "task")

var sim Simulation = Simulation{
	DstDir: filepath.Join(home, "WorkSpace", "Test"),
	SimDir: filepath.Join(home, "WorkSpace", "Test"),
	Monte:  []string{"50000"},
	Range:  Range{Start: "0.0", Step: "1.0", Stop: "2.0"},
	Signal: "n1",
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
}

func TestSetSEEDInputSPI(t *testing.T) {

	ConfigDir = filepath.Join(home, ".config", "UHA")

	p := sim.SimDir

	if err := setSEEDInputSPI(1, p, "", sim); err != nil {
		t.Fatal(err)
	}

	f, _ := ioutil.ReadDir(sim.SimDir)
	if len(f) != 1 {
		t.Fatal("Unexpected result : len(f) ", len(f))
	}

	if f[0].Name() != "50000_SEED1_input.spi" {
		t.Fatal("Unexpected result : file name : ", f[0].Name())
	}

	if err := os.Remove(filepath.Join(p, f[0].Name())); err != nil {
		t.Fatal(err)
	}
}

var nt NSeedTask = NSeedTask{
	Simulation: sim,
	Count:      3,
}

func TestSetResultDir(t *testing.T) {

	if err := setResultDir(nt, 1); err != nil {
		t.Fatal(err)
	}

	f, _ := ioutil.ReadDir(nt.Simulation.DstDir)

	if len(f) != 1 {
		t.Fatal("Unexpected result len(f) : ", len(f))
	}

	for _, v := range f {
		if !v.IsDir() || len(v.Name()) != len("RangeSEED_Sigmax.xxxx_Monte50000") {
			t.Fatal("Unexpected result")
		}
	}
}

func TestMakeSRun(t *testing.T) {
	SelfPath = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "xztaityozx", "UHA")

	actual := makeSRun(nt, 1)
	for i := 1; i <= nt.Count; i++ {

		dst := filepath.Join(nt.Simulation.DstDir, fmt.Sprintf("RangeSEED_Sigma%.4f_Monte%s/SEED%d", nt.Simulation.Vtn.Sigma, nt.Simulation.Monte[0], i))
		input := filepath.Join(nt.Simulation.SimDir, fmt.Sprintf("%s_SEED%d_input.spi", nt.Simulation.Monte[0], i))
		expect := fmt.Sprintf("cd %s && hspice -hpp -mt 4 -i %s -o ./hspice &> ./hspice.log && wv -k -ace_no_gui ../extract.ace &> wv.log && cat store.csv | sed '/^#/d;1,1d' | awk -F, '{print $2}' | xargs -n3 > ../Sigma%.4f/SEED%d.csv\n", dst, input, nt.Simulation.Vtn.Sigma, i)

		if expect != actual[i-1] {
			t.Fatal("Unexpected result : index = ", i, " : ", actual[i-1], "\nexpect : ", expect)
		}
	}
}

func TestReadNSeedTaskList(t *testing.T) {
	b, _ := json.Marshal(nt)
	p := filepath.Join(taskDir, "srun", RESERVE)

	ReserveSRunDir = p
	tryMkdir(ReserveSRunDir)

	t.Run("1", func(t *testing.T) {
		f := filepath.Join(p, "task1.json")

		log.Println(f)

		if err := ioutil.WriteFile(f, b, 0644); err != nil {
			t.Fatal(err)
		}

		actual, _ := readNSTaskFileList()

		if len(actual) != 1 {
			t.Fatal("Unexpected result : len(actual) : ", len(actual))
		}
	})

}

func TestSRun(t *testing.T) {
	srun(2, false, []NSeedTask{}, 1, []string{"echo a", "echo b", "echo c"})
}

func TestZZZ(t *testing.T) {
	f, _ := ioutil.ReadDir(nt.Simulation.DstDir)

	for _, v := range f {
		if err := os.RemoveAll(filepath.Join(nt.Simulation.DstDir, v.Name())); err != nil {
			t.Fatal(err)
		}
	}
}
