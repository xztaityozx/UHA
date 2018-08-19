package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAllConfig(t *testing.T) {
	home := os.Getenv("HOME")
	vCfg := `{"AA":"BB","CC":"DD"}`
	tCfg := `
{
  "Simulation":{
    "Monte":["100", "200", "500", "1000", "2000", "5000", "10000", "20000", "50000"],
    "Range":{
      "Start":"2.5ns",
      "Stop":"17.5ns",
      "Step":"7.5ns"
    },
    "Signal":"m8d",
    "SimDir":"/home/xztaityozx/WorkSpace/Test/",
    "DstDir":"/home/xztaityozx/WorkSpace/Test/",
    "Vtp":{
      "Voltage":-0.6,
      "Sigma":0.0,
      "Deviation":1.0
    },
    "Vtn":{
      "Voltage":0.6,
      "Sigma":0.0,
      "Deviation":1.0
    }
  },
  "TaskDir":"~/.config/UHA/task",
  "Repository":[
    {
      "Type":"Dir",
      "Path":"/home/xztaityozx/WorkSpace/Test"
    }
  ]
}
`
	vp := filepath.Join(home, "WorkSpace", "v.json")
	tp := filepath.Join(home, "WorkSpace", "t.json")

	t.Run("000_Write", func(t *testing.T) {
		if err := ioutil.WriteFile(vp, []byte(vCfg), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(tp, []byte(tCfg), 0644); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("001_Collect", func(t *testing.T) {
		if !IsAccurateConfig(tp) {
			t.Fatal("Unexpected result tp : ")
		}
	})

	t.Run("002_InVaild", func(t *testing.T) {
		if IsAccurateConfig(vp) {
			t.Fatal("Unexpected result vp : ")
		}
	})

	t.Run("00X_Remove", func(t *testing.T) {
		//os.Remove(vp)
		//os.Remove(tp)
	})
}
