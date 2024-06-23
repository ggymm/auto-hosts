package main

import (
	"bufio"
	"io"
	"net"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

var (
	hostsFile           = "hosts"
	domainsFile         = "domains.txt"
	nameserversFile     = "nameservers.txt"
	nameserversTempFile = "nameserversTemp.txt"
)

var (
	emptyList = make([]string, 0)
)

func LoadHosts() []string {
	hosts, err := readLines(hostsFile)
	if err != nil {
		log.Error().
			Str("file", hostsFile).
			Err(errors.WithStack(err)).Msg("read hosts file error")
		return emptyList
	}
	return hosts
}

func LoadDomains() []string {
	domains, err := readLines(domainsFile)
	if err != nil {
		log.Error().
			Str("file", domainsFile).
			Err(errors.WithStack(err)).Msg("read domains file error")
		return emptyList
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
		return emptyList
	}
	return nss
}

func RenewNameservers() {
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
}
