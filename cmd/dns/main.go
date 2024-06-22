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
	_ = fetchNameservers()
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
