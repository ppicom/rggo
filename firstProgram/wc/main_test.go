package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")
	exp := 4

	res := count(b, false, false)

	if res != exp {
		t.Errorf("Expected %d. got %d instead", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3\nline2\nline3 word1")
	want := 3

	got := count(b, true, false)

	if got != want {
		t.Errorf("Expected %d, got %d instead.\n", want, got)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("word1")
	want := 5

	got := count(b, false, true)

	if got != want {
		t.Errorf("Expected %d, but got %d instead.\n", want, got)
	}
}

func TestCountFromFile(t *testing.T) {
	want := 3
	file, err := os.Open("./testdata/test1.txt")
	if err != nil {
		t.Fatal(err)
	}

	got := count(file, false, false)

	if got != want {
		t.Errorf("Expected %d, but got %d instead.\n", want, got)
	}
}

func TestCountFromManyFiles(t *testing.T) {
	want := "word1 word2 word3\nword1 word2 word3\n"

	reader, err := getReader("./testdata/test1.txt,./testdata/test2.txt")
	if err != nil {
		t.Fatal(err)
	}

	got := ""
	s := bufio.NewScanner(reader)
	s.Split(bufio.ScanBytes)
	for s.Scan() {
		got += s.Text()
	}

	if got != want {
		t.Errorf("want: %q", want)
		t.Errorf("got: %q", got)
		t.Errorf("File contents are not the expected")
	}
}
