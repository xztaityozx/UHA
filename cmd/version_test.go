package cmd

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	v := UHAVersion{
		Major:    1,
		Minor:    1,
		Build:    1,
		Revision: 1,
		Status:   "Stable",
		Date:     "xxxx/xx/xx",
	}

	expect := "UHA [Ultra H_SPICE Attacker]\nVersion: 1.1.1.1 Stable (xxxx/xx/xx)\nRepository: https://github.com/xztaityozx/UHA\nLicense: MIT"
	actual := getVersion(v)

	if expect != actual {
		t.Fatal("Unexpected result : ", actual)
	}
}
