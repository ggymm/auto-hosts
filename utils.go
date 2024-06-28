package autohosts

import (
	"bufio"
	"io"
	"os"
)

func isExist(name string) bool {
	st, err := os.Stat(name)
	return !os.IsNotExist(err) && !st.IsDir()
}

func readLines(name string) ([]string, error) {
	fd, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	ls := make([]string, 0)
	buf := bufio.NewReader(fd)
	for {
		l, _, err1 := buf.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		ls = append(ls, string(l))
	}
	_ = fd.Close()
	return ls, nil
}

func writeLines(name string, lines []string) error {
	f, err := os.OpenFile(name, os.O_TRUNC|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	for _, l := range lines {
		_, _ = f.WriteString(l)
		_, _ = f.WriteString("\n")
	}
	return f.Close()
}
