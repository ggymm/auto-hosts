package main

import (
	"net"
)

func ParseGateway(ip net.IP) (net.IP, error) {
	return parse(ip)
}
