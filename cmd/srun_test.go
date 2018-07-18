package cmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestSetSEEDInputSPI(t *testing.T) {
	sim := Simulation{
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

}
