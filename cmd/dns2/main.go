package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"os"
	"slices"
)

func main() {
	name := "data/nameservers.txt"

	fd, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
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

	slices.SortFunc(src, func(i, j string) int {
		ip1 := net.ParseIP(i).To4()
		ip2 := net.ParseIP(j).To4()

		int1 := binary.BigEndian.Uint32(ip1)
		int2 := binary.BigEndian.Uint32(ip2)

		if int1 < int2 {
			return -1
		} else if int1 > int2 {
			return 1
		} else {
			return 0
		}
	})

	// 保存 nameservers.txt 文件
	fd, err = os.OpenFile(name, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for _, s := range src {
		_, _ = fd.WriteString(s + "\n")
	}
	_ = fd.Close()
}
