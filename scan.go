package autohosts

import (
	"fmt"
	"sync"
	"time"

	"github.com/ggymm/dns"
)

type Info struct {
	Rtt        time.Duration
	Addr       string
	Domain     string
	Nameserver string
}

func (i *Info) String() string {
	rtt := i.Rtt.Round(time.Millisecond).String()
	return fmt.Sprintf("%s|%s|%s", i.Addr, rtt, i.Nameserver)
}

type Scanner struct {
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (*Scanner) Scan(domain string, nameservers []string) []*Info {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	d := &sync.Map{}
	wg := &sync.WaitGroup{}
	for _, nameserver := range nameservers {
		wg.Add(1)

		go func(nameserver string) {
			defer wg.Done()

			c := new(dns.Client)
			c.Timeout = 3 * time.Second
			r, _, err := c.Exchange(m, nameserver+":53")
			if err != nil {
				return
			}
			for _, answer := range r.Answer {
				if a, ok := answer.(*dns.A); ok {
					ip := a.A.To4().String()
					if ip == "0.0.0.0" {
						continue
					}
					if _, exist := d.Load(ip); !exist {
						d.Store(ip, &Info{
							Addr:       ip,
							Domain:     domain,
							Nameserver: nameserver,
						})
					}
				}
			}
		}(nameserver)
	}
	wg.Wait()

	// 转为列表
	ips := make([]*Info, 0)
	d.Range(func(key, value any) bool {
		ips = append(ips, value.(*Info))
		return true
	})
	return ips
}
