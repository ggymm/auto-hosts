package autohosts

import (
	"fmt"
	"testing"
)

func TestScanner_Run(t *testing.T) {
	s := NewScanner()
	domains := []string{
		"github.com",
		"alive.github.com",
		"api.github.com",
	}
	nameservers, err := ReadLines("data/nameservers.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, domain := range domains {
		ips := s.Scan(domain, nameservers)
		for _, ip := range ips {
			fmt.Println(ip)
		}
	}
}
