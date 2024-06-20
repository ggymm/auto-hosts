package main

import (
	"net"
	"testing"

	"auto-hosts/log"
)

func TestScanner_NetInfo(t *testing.T) {
	log.Init()

	d := Device{
		IpAddr: net.IP{192, 168, 1, 27},
	}
	mac, err := net.ParseMAC("54:05:db:83:7f:a5")
	if err != nil {
		t.Fatal(err)
	}
	d.HwAddr = mac
	d.Name = `\Device\NPF_{81A86FFA-2C4F-4E6B-AD4E-29036647FB75}`
	d.Desc = "Realtek PCIe GbE Family Controller"

	s := NewScanner()
	err = s.Init(&d)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("gwIp: %s, gwHw: %s", s.gwIp, s.gwHw)
}
