package main

import (
	"net"
	"testing"
)

func TestParseGateway(t *testing.T) {
	t.Log(ParseGateway(net.IP{192, 168, 1, 102}))
}
