package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAllCount(t *testing.T) {
	home := os.Getenv("HOME")
	file := filepath.Join(home, "WorkSpace", "out.csv")

	os.MkdirAll(filepath.Join(home, "WorkSpace"), 0777)
	ioutil.WriteFile(file, []byte("0.8 0.0 0.0\n0.8 0.8 0.8"), 0644)

	t.Run("001_NewAggregateData", func(t *testing.T) {
		ad, err := NewAggregateData(file)
		if err != nil {
			t.Fatal(err)
		}
		if ad.Failure != 1 {
			t.Fatal("Fatal : ", ad.Failure)
		}
		if ad.Lines != 2 {
			t.Fatal("Fatal : ", ad.Lines)
		}
	})

}
