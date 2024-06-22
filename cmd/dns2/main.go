package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ggymm/dns"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

var (
	nameserversFile = "nameservers.txt"
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
	fd, err := os.OpenFile(nameserversFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
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
				r, _, err1 := c.Exchange(m, s+":53")
				if err1 != nil {
					if strings.Contains(err1.Error(), "timeout") {
						return
					}
					fmt.Println(err1)
					return
				}
				dst = append(dst, s)
				if len(r.Answer) != 0 {
					for _, rr := range r.Answer {
						fmt.Println(rr.String())
					}
				}
			}(s)
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
	}

	// 保存 nameservers.txt 文件
	f1, err := os.OpenFile(nameserversFile, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("open nameservers file error")
		return
	}
	for i, s := range dst {
		_, _ = f1.WriteString(s + "\n")
		log.Info().Msgf("%d: %s", i, s)
	}
}
