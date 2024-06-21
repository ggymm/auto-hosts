package main

import (
	"fmt"

	"github.com/ggymm/ping"
)

func main() {
	pinger, err := ping.NewPinger("1.1.1.1")
	pinger.SetPrivileged(true)
	if err != nil {
		panic(err)
	}
	pinger.Count = 4
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics()
	fmt.Println(stats)
}
