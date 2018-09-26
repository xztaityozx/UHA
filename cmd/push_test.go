package cmd

import "testing"

func TestPushAll(t *testing.T) {
	t.Run("netxColumn A->B", func(t *testing.T) {
		expect := "B"
		actual := NextColumn("A")

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})
	t.Run("netxColumn Z->AA", func(t *testing.T) {
		expect := "AA"
		actual := NextColumn("Z")

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})
	t.Run("netxColumn AD->AE", func(t *testing.T) {
		expect := "AE"
		actual := NextColumn("AD")

		if actual != expect {
			t.Fatal(actual, "is not", expect)
		}
	})
}
