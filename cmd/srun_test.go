package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var sim Simulation = Simulation{
	DstDir: "/home/xztaityozx/WorkSpace/Test/",
	SimDir: "/home/xztaityozx/WorkSpace/Test/",
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

	ConfigDir = "/home/xztaityozx/.config/UHA"

	p := filepath.Join(sim.SimDir, "50000_SEED1_input.spi")

	if err := setSEEDInputSPI(1, p, sim); err != nil {
		t.Fatal(err)
	}

	f, _ := ioutil.ReadDir(sim.SimDir)
	if len(f) != 1 {
		t.Fatal("Unexpected result : len(f) ", len(f))
	}

	if f[0].Name() != "50000_SEED1_input.spi" {
		t.Fatal("Unexpected result : file name : ", f[0].Name())
	}

	os.Remove(p)
}

var nt NSeedTask = NSeedTask{
	Simulation: sim,
	Count:      3,
}

func TestSetResultDir(t *testing.T) {

	if err := setResultDir(nt); err != nil {
		t.Fatal(err)
	}

	f, _ := ioutil.ReadDir(nt.Simulation.DstDir)

	if len(f) != 3 {
		t.Fatal("Unexpected result len(f) : ", len(f))
	}

	for _, v := range f {
		if !v.IsDir() || len(v.Name()) != len("Monte50000_SEEDx") {
			t.Fatal("Unexpected result")
		}
	}
}

func TestMakeSRun(t *testing.T) {
	expect := "cd /home/xztaityozx/WorkSpace/Test"
	actual := makeSRun(nt)
}
