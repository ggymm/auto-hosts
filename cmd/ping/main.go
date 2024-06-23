package main

import (
	"fmt"
	"time"

	"github.com/ggymm/ping"
)

func main() {
	pinger, err := ping.NewPinger("61.216.168.145")
	pinger.Count = 4
	pinger.Timeout = 1 * time.Second
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics()
	fmt.Printf("%+v\n", stats)
}
