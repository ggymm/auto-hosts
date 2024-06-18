package main

import (
	"net"

	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

type Device struct {
	Name string // pcap 设备名称
	Desc string // pcap 设备描述

	IpAddr net.IP
	HwAddr net.HardwareAddr
}

func (d *Device) String() string {
	return d.Desc + " - " + d.IpAddr.String()
}

func GetDevices() []*Device {
	var (
		err  error
		ift  []net.Interface
		devs []pcap.Interface
	)
	ret := make([]*Device, 0)
	ift, err = net.Interfaces()
	if err != nil {
		log.Error().Err(errors.WithStack(err)).Msg("find interface error")
		return ret
	}
	for _, i := range ift {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if n, ok := addr.(*net.IPNet); !ok {
				continue
			} else {
				if n.IP.To4() != nil && !n.IP.IsLoopback() {
					ret = append(ret, &Device{
						IpAddr: n.IP.To4(),
						HwAddr: i.HardwareAddr,
					})
					break
				}
			}
		}
	}

	devs, err = pcap.FindAllDevs()
	if err != nil {
		log.Error().Err(errors.WithStack(err)).Msg("find pcap devices error")
		return ret
	}
	for _, dev := range devs {
		ip := net.IP{}
		for _, addr := range dev.Addresses {
			if addr.IP.To4() != nil {
				ip = addr.IP.To4()
				break
			}
		}
		for _, r := range ret {
			if r.IpAddr.Equal(ip) {
				r.Name = dev.Name
				r.Desc = dev.Description
				break
			}
		}
	}
	return ret
}
