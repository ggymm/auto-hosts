package main

import (
	"auto-hosts/pkg"
	"slices"
	"sync"
	"time"

	"github.com/ggymm/dns"
)

type Scanner struct {
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (*Scanner) Scan(domains, nameservers []string) map[string][]string {
	ret := make(map[string][]string)
	for _, domain := range domains {
		ret[domain] = make([]string, 0)
	}

	for _, domain := range domains {
		wg := &sync.WaitGroup{}
		ips := pkg.Slice[string]{}
		for _, nameserver := range nameservers {
			wg.Add(1)

			go func(domain, nameserver string) {
				defer wg.Done()

				m := new(dns.Msg)
				m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
				m.RecursionDesired = true

				c := new(dns.Client)
				c.Timeout = 1 * time.Second
				r, _, err := c.Exchange(m, nameserver+":53")
				if err != nil {
					return
				}
				for _, answer := range r.Answer {
					if a, ok := answer.(*dns.A); ok {
						ips.Append(a.A.To4().String())
					}
				}
			}(domain, nameserver)
		}
		wg.Wait()
		time.Sleep(1 * time.Second)

		// 收集结果
		newIps := make([]string, 0)
		ips.Foreach(func(i int, ip string) {
			if !slices.Contains(newIps, ip) {
				newIps = append(newIps, ip)
			}
		})
		ret[domain] = newIps
	}
	return ret
}
