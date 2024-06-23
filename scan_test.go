package main

import (
	"auto-hosts/log"
	"fmt"
	"testing"
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

	ret := s.Scan(domains, nameservers)
	for d, ips := range ret {
		fmt.Println(d)
		for _, ip := range ips {
			fmt.Println(ip)
		}
	}
}
