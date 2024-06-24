package main

import (
	"fmt"
	"time"

	"github.com/ggymm/ping"
)

func main() {
	p, err := ping.NewPinger("149.126.86.55")
	p.Count = 4
	p.Timeout = 1 * time.Second
	p.SetPrivileged(true)
	err = p.Run()
	if err != nil {
		panic(err)
	}
	stats := p.Statistics()
	if stats.PacketsRecv == 0 {
		fmt.Println("Destination unreachable")
		return
	}
	fmt.Println(stats.AvgRtt.String())
}
