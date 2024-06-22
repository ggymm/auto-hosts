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
	m.SetQuestion(dns.Fqdn("github.com"), dns.TypeA)
	m.RecursionDesired = true

	wg := &sync.WaitGroup{}
	dst := make([]string, 0)
	for _, s := range src {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()

			c := new(dns.Client)
			c.Timeout = 1 * time.Second
			r, _, err1 := c.Exchange(m, s+":53")
			if err1 != nil {
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

	for i, s := range dst {
		fmt.Println(i, s)
	}
}
