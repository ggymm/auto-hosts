package main

import (
	"slices"
	"sync"
	"time"

	"github.com/ggymm/dns"
	"github.com/pkg/errors"

	"auto-hosts/log"
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

		es := make([]string, 0)
		ips := make([]string, 0)
		for _, nameserver := range nameservers {
			wg.Add(1)

			go func(domain, nameserver string) {
				defer wg.Done()

				m := new(dns.Msg)
				m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
				m.RecursionDesired = true

				c := new(dns.Client)
				c.Timeout = 3 * time.Second
				r, _, err := c.Exchange(m, nameserver+":53")
				if err != nil {
					log.Error().
						Str("domain", domain).
						Str("nameserver", nameserver).
						Err(errors.WithStack(err)).Msg("dns query failed")

					// 添加到异常列表
					es = append(es, nameserver)
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

		// 移除异常
		for _, e := range es {
			for i, nameserver := range nameservers {
				if e == nameserver {
					// 删除当前元素
					nameservers = append(
						nameservers[:i],
						nameservers[i+1:]...,
					)
				}
			}
		}

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
