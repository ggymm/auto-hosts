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
	for k, _ := range m {
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
