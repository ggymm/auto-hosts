package main

import (
	"bufio"
	"io"
	"os"
	"slices"
	"strings"
)

func main() {
	name := "data/domains.txt"

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
		if len(l) == 0 {
			continue
		}
		src = append(src, string(l))
	}
	_ = fd.Close()

	slices.SortFunc(src, func(i, j string) int {
		println(i, j)
		return strings.Compare(i, j)
	})

	fd, err = os.OpenFile(name, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for _, s := range src {
		_, _ = fd.WriteString(s + "\n")
	}
	_ = fd.Close()
}
