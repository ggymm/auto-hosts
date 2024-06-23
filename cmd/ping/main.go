package main

import (
	"fmt"
	"time"

	"github.com/ggymm/ping"
)

func main() {
	p, err := ping.NewPinger("8.8.8.8")
	p.Count = 4
	p.Timeout = 1 * time.Second
	p.SetPrivileged(true)
	err = p.Run()
	if err != nil {
		panic(err)
	}
	stats := p.Statistics()
	fmt.Println(stats.AvgRtt.String())
}
