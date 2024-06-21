package main

import (
	"os"
	"testing"
)

func Test_Data_GetNss(t *testing.T) {
	err := os.Chdir("data")
	if err != nil {
		t.Fatal(err)
	}

	nss := LoadNameservers()
	for i, ns := range nss {
		t.Logf("%d: %s", i, ns)
	}
}

func Test_Data_GetDomains(t *testing.T) {
	err := os.Chdir("data")
	if err != nil {
		t.Fatal(err)
	}

	domains := LoadDomains()
	for i, domain := range domains {
		t.Logf("%d: %s", i, domain)
	}
}
