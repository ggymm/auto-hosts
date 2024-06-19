package main

import (
	"testing"
)

func Test_GetDevices(t *testing.T) {
	devices := GetDevices()
	for _, d := range devices {
		t.Logf("%+v", d)
	}
}
