package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ggymm/ping"
)

var (
	ips = []string{
		"20.26.156.215",
		"140.82.116.3",
		"140.82.112.3",
		"140.82.121.4",
		"20.27.177.113",
		"20.205.243.166",
		"20.200.245.247",
		"140.82.121.3",
	}
	infos = make([]*Info, 0)
)

type Info struct {
	ip  string
	rtt time.Duration
}

func main() {
	for _, ip := range ips {
		infos = append(infos, &Info{
			ip:  ip,
			rtt: 0,
		})
	}

	//fmt.Println("test")
	//test()
	//time.Sleep(1 * time.Second)

	fmt.Println("test1")
	test1()

	for _, info := range infos {
		fmt.Printf("IP: %s, RTT: %v\n", info.ip, info.rtt)
	}
}

func test() {
	for _, ip := range ips {
		p, err := ping.NewPinger(ip)
		p.Count = 4
		p.Timeout = 1 * time.Second
		p.SetPrivileged(true)
		err = p.Run()
		if err != nil {
			continue
		}
		stats := p.Statistics()
		if stats.PacketsRecv == 0 {
			continue
		}
		fmt.Printf("IP: %s, RTT: %v\n", ip, stats.AvgRtt)
	}
}

func test1() {
	wg := &sync.WaitGroup{}
	for _, info := range infos {
		wg.Add(1)
		fmt.Println(info.ip)

		go func(info *Info) {
			defer wg.Done()

			p, err := ping.NewPinger(info.ip)
			p.Count = 4
			p.Timeout = 1 * time.Second
			p.SetPrivileged(true)
			err = p.Run()
			if err != nil {
				return
			}
			stats := p.Statistics()
			if stats.PacketsRecv == 0 {
				info.rtt = 99 * time.Second
			} else {
				info.rtt = stats.AvgRtt
			}
		}(info)
	}
	wg.Wait()
}
