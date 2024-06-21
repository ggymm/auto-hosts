package main

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
		"assets-cdn.github.com",
		"avatars.githubusercontent.com",
		"avatars0.githubusercontent.com",
		"avatars1.githubusercontent.com",
		"avatars2.githubusercontent.com",
		"avatars3.githubusercontent.com",
		"avatars4.githubusercontent.com",
		"avatars5.githubusercontent.com",
		"camo.githubusercontent.com",
	}
	nameservers, err := readLines("data/nameservers.txt")
	if err != nil {
		t.Fatal(err)
	}

	ret := s.Run(domains, nameservers)
	for d, ips := range ret {
		fmt.Println(d)
		for _, ip := range ips {
			fmt.Println(ip)
		}
	}
}
