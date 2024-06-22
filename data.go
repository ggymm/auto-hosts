package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ggymm/dns"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

var (
	domainsFile         = "domains.txt"
	nameserversFile     = "nameservers.txt"
	nameserversTempFile = "nameserversTemp.txt"
)

func LoadDomains() []string {
	domains, err := readLines(domainsFile)
	if err != nil {
		return nil
	}
	return domains
}

func LoadNameservers() []string {
	// 读取文件
	nss, err := readLines(nameserversFile)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("read nameservers file error")
		return nss
	}
	return nss
}

func GetNameservers() {
	FetchNameservers()
	FilterNameservers()
}

func FetchNameservers() {
	url := "https://public-dns.info/nameservers.txt"
	_, err := resty.New().R().SetOutput(nameserversFile).Get(url)
	if err != nil {
		log.Error().
			Str("url", url).
			Err(errors.WithStack(err)).Msg("download nameservers error")
		return
	}

	var (
		f1 *os.File
		f2 *os.File
	)
	f1, err = os.OpenFile(nameserversFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("read nameservers file error")
		return
	}
	f2, err = os.OpenFile(nameserversTempFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversTempFile).
			Err(errors.WithStack(err)).Msg("create nameservers file error")
	}
	buf := bufio.NewReader(f1)
	for {
		l, _, err1 := buf.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		ip := net.ParseIP(string(l))
		if ip.To4() != nil {
			_, _ = f2.Write(l)
			_, _ = f2.WriteString("\n")
		}
	}
	_ = f1.Close()
	_ = f2.Close()

	// 重命名文件
	err = os.Rename(nameserversTempFile, nameserversFile)
	if err != nil {
		log.Error().
			Str("file1", nameserversFile).
			Str("file2", nameserversTempFile).
			Err(errors.WithStack(err)).Msg("rename nameservers file error")
		return
	}
	return
}

func FilterNameservers() {
	fd, err := os.OpenFile(nameserversFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("read nameservers file error")
		return
	}
	src := make([]string, 0)
	buf := bufio.NewReader(fd)
	for {
		l, _, err1 := buf.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		src = append(src, string(l))
	}
	_ = fd.Close()

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	m.RecursionDesired = true

	dst := make([]string, 0)
	size := 10000
	for i := 0; i < len(src); i += size {
		end := i + size
		if end > len(src) {
			end = len(src)
		}

		wg := &sync.WaitGroup{}
		group := src[i:end]
		for _, s := range group {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()

				c := new(dns.Client)
				c.Timeout = 1 * time.Second
				_, _, err1 := c.Exchange(m, s+":53")
				if err1 != nil {
					if strings.Contains(err1.Error(), "timeout") {
						return
					}
					fmt.Println(err1)
					return
				}
				dst = append(dst, s)
			}(s)
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
	}

	// 保存 nameservers.txt 文件
	fd, err = os.OpenFile(nameserversFile, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("open nameservers file error")
		return
	}
	for _, s := range dst {
		_, _ = fd.WriteString(s + "\n")
	}
	_ = fd.Close()
}
