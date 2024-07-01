package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ggymm/ping"
	"github.com/pkg/errors"

	"auto-hosts"
)

var (
	scanner *autohosts.Scanner

	domains     []string
	nameservers []string
)

func init() {
	dir := ""
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Base(exe)
	if !strings.HasPrefix(exe, os.TempDir()) && !strings.HasPrefix(path, "___") {
		dir = filepath.Dir(exe)
	} else {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			// 需要根据当前文件所处目录，修改相对位置
			dir = filepath.Join(filepath.Dir(filename), "../..")
		}
	}

	// 设置 app 工作目录
	err = os.Chdir(filepath.Join(dir, "data"))
	if err != nil {
		panic(errors.WithStack(err))
	}

	scanner = autohosts.NewScanner()
}

func main() {
	// 读取信息
	domains = autohosts.LoadDomains()
	nameservers = autohosts.LoadNameservers()

	// 启动程序
	start()
}

func start() {
	fmt.Println("start")

	hosts := make([]string, 0)
	for i, domain := range domains {
		fmt.Printf("scan %d/%d: %s\n", i+1, len(domains), domain)

		ip := ""
		list := scanner.Scan(domain, nameservers)
		if len(list) > 0 {
			wg := &sync.WaitGroup{}
			for _, item := range list {
				wg.Add(1)
				fmt.Printf("ping %s\n", item.Addr)

				go func(item *autohosts.Info) {
					defer wg.Done()

					p, _ := ping.NewPinger(item.Addr)
					p.Count = 4
					p.Timeout = 1 * time.Second
					p.SetPrivileged(true)
					err := p.Run()
					if err != nil {
						return
					}
					stats := p.Statistics()
					if stats.PacketsRecv != 0 {
						item.Rtt = stats.AvgRtt
					} else {
						item.Rtt = 99 * time.Second
					}
				}(item)
			}
			wg.Wait()

			// 排序
			slices.SortFunc(list, func(i, j *autohosts.Info) int {
				if i.Rtt < j.Rtt {
					return -1
				} else {
					return 1
				}
			})
			ip = list[0].Addr + " " + domain

			// 保存到文件
			ips := make([]string, 0)
			for _, item := range list {
				ips = append(ips, item.String())
				fmt.Printf("ping rtt: %s\n", item.String())
			}
			err := autohosts.WriteLines(fmt.Sprintf("ips/%s.txt", domain), ips)
			if err != nil {
				panic(err)
			}
		} else {
			ip = "unknown"
		}
		hosts = append(hosts, ip)
	}

	// 打印结果
	fmt.Println("finish")
	for _, host := range hosts {
		fmt.Printf("%s\n", host)
	}

	// 保存到文件
	err := autohosts.WriteLines("hosts", hosts)
	if err != nil {
		panic(err)
	}
}
