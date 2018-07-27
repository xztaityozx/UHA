package cmd

import (
	"testing"
)

func TestGetMultiLine(t *testing.T) {
	li := []listItem{
		listItem{
			Date: "1111/22/33/44:55",
			Name: "1",
		},
		listItem{
			Date: "1111/22/33/44:55",
			Name: "2",
		},
		listItem{
			Date: "1111/22/33/44:55",
			Name: "3",
		},
	}

	actual := getMultiLine(li, 1)
	expect := []string{
		"\t1",
		"\t2",
		"\t3",
	}

	for i := 0; i < 3; i++ {
		if actual[i] != expect[i] {
			t.Fatal("Unexpected result : ", actual[i], "\nexpect : ", expect[i])
		}
	}
}

func TestGetLongList(t *testing.T) {
	li := []listItem{
		listItem{
			Date: "1111/22/33/44:55",
			Name: "20180727153226_1",
		},
		listItem{
			Date: "2222/22/33/44:55",
			Name: "20180727153226_2",
		},
		listItem{
			Date: "3333/22/33/44:55",
			Name: "20180727153226_3",
		},
	}
	actual := getLongList(li)
	expect := []string{
		"1\t1111/22/33/44:55",
		"2\t2222/22/33/44:55",
		"3\t3333/22/33/44:55",
	}
	for i := 0; i < 3; i++ {
		if actual[i] != expect[i] {
			t.Fatal("Unexpected result : ", actual[i], "\nexpect : ", expect[i])
		}
	}
}

func TestGetSingleLineList(t *testing.T) {
	li := []listItem{
		listItem{
			Date: "1111/22/33/44:55",
			Name: "1",
		},
		listItem{
			Date: "2222/22/33/44:55",
			Name: "2",
		},
		listItem{
			Date: "3333/22/33/44:55",
			Name: "3",
		},
	}

	actual := getSingleLineList(li)
	expect := []string{
		"1",
		"2",
		"3",
	}
	for i := 0; i < 3; i++ {
		if actual[i] != expect[i] {
			t.Fatal("Unexpected result : ", actual[i], "\nexpect : ", expect[i])
		}
	}

}
