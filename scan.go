package main

import (
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

func (*Scanner) Run(domains, nameservers []string) map[string][]string {
	ret := make(map[string][]string)
	for _, domain := range domains {
		ret[domain] = make([]string, 0)
	}

	for _, domain := range domains {
		wg := &sync.WaitGroup{}
		time.Sleep(1 * time.Second)

		ips := make([]string, 0)
		for _, nameserver := range nameservers {
			wg.Add(1)

			go func(domain, nameserver string) {
				defer wg.Done()

				m := new(dns.Msg)
				m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
				m.RecursionDesired = true

				c := new(dns.Client)
				c.Timeout = 5 * time.Second
				r, _, err := c.Exchange(m, nameserver+":53")
				if err != nil {
					return
				}
				if len(r.Answer) == 0 {
					return
				}
				for _, answer := range r.Answer {
					if a, ok := answer.(*dns.A); ok {
						ips = append(ips, a.A.To4().String())
					}
				}
			}(domain, nameserver)
		}
		wg.Wait()

		// 收集结果
		newIps := make([]string, 0)
		for _, ip := range ips {
			if !slices.Contains(newIps, ip) {
				newIps = append(newIps, ip)
			}
		}
		ret[domain] = newIps
	}
	return ret
}
