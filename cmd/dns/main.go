package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

var (
	nameserversFile     = "nameservers.txt"
	nameserversTempFile = "nameserversTemp.txt"
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
			dir = filepath.Join(filepath.Dir(filename), "../../")
		}
	}
	wd := filepath.Join(dir, "temp")

	// 设置 app 工作目录
	err = os.Chdir(wd)
	if err != nil {
		panic(errors.WithStack(err))
	}

	log.Init()
}

func main() {
	// 删除 nameservers.txt 文件
	_ = os.Remove(nameserversFile)

	// 获取最新的 nameservers 并且按照响应时间排序
	nss := fetchNameservers()
	//nss = filterNameservers(nss)

	// 保存 nameservers.txt 文件
	//f1, err := os.OpenFile(nameserversFile, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	//if err != nil {
	//	log.Error().
	//		Str("file", nameserversFile).
	//		Err(errors.WithStack(err)).Msg("open nameservers file error")
	//	return
	//}
	for i, s := range nss {
		//_, _ = f1.WriteString(s + "\n")
		log.Info().Msgf("%d: %s", i, s)
	}
}

func fetchNameservers() (nss []string) {
	nss = make([]string, 0)
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
	fb := bufio.NewReader(f1)
	for {
		l, _, err1 := fb.ReadLine()
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

			nss = append(nss, string(l))
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

//func filterNameservers(src []string) (dst []string) {
//	dst = make([]string, 0)
//
//	wg := new(sync.WaitGroup)
//	for _, ns := range src {
//		wg.Add(1)
//		go func(ns string) {
//			defer wg.Done()
//
//			p, err := ping.NewPinger(ns)
//			if err != nil {
//				log.Error().
//					Str("remote", ns).
//					Err(errors.WithStack(err)).Msg("ping nameserver error")
//				return
//			}
//			p.SetPrivileged(true)
//			p.Count = 1
//			p.Timeout = 1 * time.Second
//			err = p.Run()
//			if err != nil {
//				log.Error().
//					Str("remote", ns).
//					Err(errors.WithStack(err)).Msg("ping nameserver error")
//				return
//			}
//			dst = append(dst, ns)
//		}(ns)
//	}
//	wg.Wait()
//	return dst
//}

//func filterNameservers2(src []string) (dst []string) {
//	dst = make([]string, 0)
//
//	c := new(dns.Client)
//	c.Timeout = 1 * time.Second
//	for _, ns := range src {
//
//		m := new(dns.Msg)
//		m.SetQuestion("google.com.", dns.TypeA)
//		m.RecursionDesired = true
//
//		// 发送 DNS 查询并接收响应
//		_, _, err := c.Exchange(m, ns+":53")
//		if err != nil {
//			log.Error().
//				Str("remote", ns).
//				Err(errors.WithStack(err)).Msg("nameservers not available")
//			continue
//		}
//		dst = append(dst, ns)
//	}
//	return dst
//}
