package main

import (
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
