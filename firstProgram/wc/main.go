package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	fname := flag.String("f", "", "Comma separated files to read from")
	flag.Parse()

	reader, err := getReader(*fname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(count(reader, *lines, *bytes))
}

func getReader(fnames string) (io.Reader, error) {
	if fnames == "" {
		return os.Stdin, nil
	}

	var data []byte

	files := strings.Split(fnames, ",")

	for _, file := range files {
		d, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		data = append(data, d...)
	}

	return bytes.NewReader(data), nil
}

func count(r io.Reader, countLines, countBytes bool) int {
	scanner := bufio.NewScanner(r)

	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	if countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	wc := 0

	for scanner.Scan() {
		wc += 1
	}

	return wc
}
