package main

import (
	"time"

	"github.com/ggymm/dns"
)

func main() {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	m.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 3 * time.Second
	r, _, err := c.Exchange(m, "61.216.168.145:53")
	if err != nil {
		panic(err)
	}
	for _, answer := range r.Answer {
		println(answer.String())
	}
}
