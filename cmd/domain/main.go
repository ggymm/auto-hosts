package main

import (
	"time"

	"github.com/ggymm/dns"
)

func main() {
	m := new(dns.Msg)
	m.SetQuestion("github.com.", dns.TypeA)
	m.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 3 * time.Second
	r, _, err := c.Exchange(m, "45.76.64.64:53")
	if err != nil {
		panic(err)
	}
	if r.Rcode != dns.RcodeSuccess {
		panic(err)
	}
	for _, answer := range r.Answer {
		println(answer.String())
	}
}
