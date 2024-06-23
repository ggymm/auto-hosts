package main

import (
	"fmt"
	"testing"

	"auto-hosts/log"
)

func TestScanner_Run(t *testing.T) {
	log.Init("temp/app")

	s := NewScanner()
	domains := []string{
		"github.com",
		"alive.github.com",
		"api.github.com",
	}
	nameservers, err := readLines("data/nameservers.txt")
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
