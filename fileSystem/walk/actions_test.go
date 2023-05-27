package main

import (
	"os"
	"testing"
	"time"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tc.file, tc.ext, tc.minSize, info, "")

			if f != tc.expected {
				t.Errorf("Expected '%t', got '%t' instead\n", tc.expected, f)
			}
		})
	}
}

func TestFilterOut_Mod(t *testing.T) {

	testCases := []struct {
		name     string
		file     string
		mod      string
		expected bool
	}{
		{"FilterNoModificationDate", "testdata/dir.log", "", false},
		{"FilterModificationDateNoMatch", "testdata/dir.log", "1990-01-01", true},
		{"FilterModificationDateMatch", "testdata/dir.log", time.Now().Format(time.DateOnly), false},
	}

	for _, tc := range testCases {
		t.Run(t.Name(), func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tc.file, "", 0, info, tc.mod)

			if f != tc.expected {
				t.Fatalf("Expecter '%t', got '%t' instead\n", tc.expected, f)
			}
		})
	}
}
