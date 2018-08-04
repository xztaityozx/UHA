package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAllMake(t *testing.T) {
	home := os.Getenv("HOME")
	ReserveRunDir = filepath.Join(home, ".config", "UHA", "task", "run", RESERVE)
	var sim Simulation = Simulation{
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
	}
	t.Run("001_writeTask", func(t *testing.T) {
		res, err := writeTask(sim)
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadFile(res)
		if err != nil {
			t.Fatal(err)
		}

		var s Simulation
		if err := json.Unmarshal(b, &s); err != nil {
			t.Fatal(err)
		}

		if s.SimDir != sim.SimDir {
			t.Fatal("Unexpected Result SimDir : ", s.SimDir)
		}

		os.Remove(res)

	})
}
