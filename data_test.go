package autohosts

import (
	"os"
	"testing"
)

func Test_Data_LoadDomains(t *testing.T) {
	err := os.Chdir("data")
	if err != nil {
		t.Fatal(err)
	}

	ss := LoadDomains()
	for i, s := range ss {
		t.Logf("%d: %s", i, s)
	}
}

func Test_Data_LoadNameservers(t *testing.T) {
	err := os.Chdir("data")
	if err != nil {
		t.Fatal(err)
	}

	ss := LoadNameservers()
	for i, s := range ss {
		t.Logf("%d: %s", i, s)
	}
}

func Test_Data_RenewNameservers(t *testing.T) {
	err := os.Chdir("data")
	if err != nil {
		t.Fatal(err)
	}

	RenewNameservers()
	ss := LoadNameservers()
	for i, s := range ss {
		t.Logf("%d: %s", i, s)
	}
}
