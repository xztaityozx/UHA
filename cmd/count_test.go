package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestAllCount(t *testing.T) {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, "WorkSpace")
	file := filepath.Join(home, "WorkSpace", "out.csv")

	os.MkdirAll(filepath.Join(home, "WorkSpace"), 0777)
	text := []byte("0.8 0.0 0.0\n0.8 0.8 0.8")
	ioutil.WriteFile(file, text, 0644)

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
		os.Remove(file)
	})

	t.Run("002_GetAggregateDataAll", func(t *testing.T) {
		for _, v := range []string{"SEED001.csv", "SEED002.csv", "SEED003.csv"} {
			p := filepath.Join(dir, v)
			ioutil.WriteFile(p, text, 0644)
		}

		actual := GetAggregateDataAll(dir)
		if len(actual) != 3 {
			t.Fatal("Fatal : len(actual)", len(actual))
		}
		var expects []AggregateData = []AggregateData{
			AggregateData{
				Failure:  1,
				Lines:    2,
				FileName: "SEED001.csv",
			},
			AggregateData{
				Failure:  1,
				Lines:    2,
				FileName: "SEED002.csv",
			},
			AggregateData{
				Failure:  1,
				Lines:    2,
				FileName: "SEED003.csv",
			},
		}

		log.Println(actual)

		for i, e := range expects {
			if !e.Compare(actual[i]) {
				t.Fatal(e, "is not", actual[i])
			}
		}
	})

	t.Run("003_CumulativeSum", func(t *testing.T) {
		ads := GetAggregateDataAll(dir)
		actual := CumulativeSum(&ads)
		for i, v := range []AggregateData{
			AggregateData{
				Failure: 1,
				Lines:   2,
			},
			AggregateData{
				Failure: 2,
				Lines:   4,
			},
			AggregateData{
				Failure: 3,
				Lines:   6,
			},
		} {
			if !v.Compare(actual[i]) {
				t.Fatal(v, "is not", actual[i])
			}
		}
	})

	t.Run("004_GetSigma", func(t *testing.T) {
		RangeSEEDCount = true
		actual := GetSigma("/path/to/RangeSEED_Vtn0.0000Vtp0.0000_Sigma1.1111_Monte0000/Result")
		expect := 1.1111
		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
		RangeSEEDCount = false
		actual = GetSigma("/path/to/VtpVolt0.000_VtnVolt0.000/Sigma1.1111/SEED000/Result")
		expect = 1.1111
		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}

	})
}
