package main

import (
	"encoding/binary"
	"net"
	"slices"
	"testing"
	"time"
)

func Test_Compare(t *testing.T) {
	m := map[string]time.Duration{
		"1s": time.Second,
		"2s": 2 * time.Second,
		"3s": 3 * time.Second,
	}
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}

	slices.SortFunc(l, func(i, j string) int {
		a, _ := m[i]
		b, _ := m[j]

		if a < b {
			return -1
		} else {
			return 1
		}
	})

	for _, s := range l {
		t.Log(s)
	}
}

func Test_Slice(t *testing.T) {
	list := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}
	for n := 0; n < 2; n++ {
		e := make([]string, 0)
		for i, s := range list {
			t.Logf("%d, %s", i, s)
			if s == "2" {
				e = append(e, s)
			}
			if s == "3" {
				e = append(e, s)
			}
		}
		for _, s := range e {
			for i, s1 := range list {
				if s == s1 {
					list = append(list[:i], list[i+1:]...)
				}
			}
		}
		t.Log()
	}
}

func Test_Sort(t *testing.T) {
	list := []string{
		"1.1.1.1",
		"1.2.4.8",
		"8.8.8.8",
		"185.222.222.222",
		"45.11.45.11",
		"101.101.101.101",
		"94.140.14.14",
		"223.5.5.5",
		"119.29.29.29",
		"180.76.76.76",
		"101.226.4.6",
		"114.114.114.114",
		"208.67.222.222",
		"9.9.9.9",
	}
	slices.SortFunc(list, func(i, j string) int {
		ip1 := net.ParseIP(i).To4()
		ip2 := net.ParseIP(j).To4()

		int1 := binary.BigEndian.Uint32(ip1)
		int2 := binary.BigEndian.Uint32(ip2)

		if int1 < int2 {
			return -1
		} else if int1 > int2 {
			return 1
		} else {
			return 0
		}
	})
	for _, s := range list {
		t.Log(s)
	}
}
