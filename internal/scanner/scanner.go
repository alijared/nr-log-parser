package scanner

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"time"
)

func New(f *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	return scanner
}

func ParseTime(line []byte) (time.Time, error) {
	fields := bytes.Fields(line)
	t, err := time.Parse(time.RFC3339, getStringInBetween(string(fields[0]), "\"", "\""))
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func GetLogLevel(line []byte) string {
	return getStringInBetween(string(line), "level=", " ")
}

func getStringInBetween(str string, start string, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	return str[s : s+e]
}
