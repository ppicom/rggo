package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo, mod string) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if mod != "" {
		t, err := time.Parse(time.DateOnly, mod)
		if err == nil && info.ModTime().After(t) {
			return true
		}
	}

	if ext == "" {
		return false
	}

	for _, e := range strings.Split(ext, ",") {
		if e != "" && filepath.Ext(path) == e {
			return false
		}
	}

	return true
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zw := gzip.NewWriter(out)

	zw.Name = filepath.Base(path)

	if _, err := io.Copy(zw, in); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return out.Close()
}

func restore(archiveDir, root, path string) error {
	stats, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !stats.IsDir() {
		return fmt.Errorf("%s is not a directory", root)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	if ext != ".gz" {
		return nil
	}

	dest := strings.TrimSuffix(path, ext)
	targetPath := filepath.Join(relDir, dest)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zw, err := gzip.NewReader(in)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, zw); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return out.Close()
}
