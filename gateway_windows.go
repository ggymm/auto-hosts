package main

import (
	"net"
	"os/exec"
	"strings"
	"syscall"
)

func parse(ip net.IP) (net.IP, error) {
	cmd := exec.Command("route", "print", "0.0.0.0")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	n := 0
	dst := net.IP{}
	lines := strings.Split(string(out), "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "======") {
			n++
			continue
		}
		if n == 3 {
			if len(lines) <= i+2 {
				break
			}
			l = lines[i+2]
			if strings.HasPrefix(l, "======") {
				break
			}
			fields := strings.Fields(l)
			if len(fields) >= 5 &&
				fields[0] == "0.0.0.0" &&
				fields[3] == ip.String() {
				dst = net.ParseIP(fields[2]).To4()
				break
			}
		}
	}
	return dst, nil
}
