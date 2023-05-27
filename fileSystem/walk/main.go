package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext     string
	archive string
	mod     string
	size    int64
	list    bool
	del     bool
	restore bool
	wLog    io.Writer
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	logFile := flag.String("log", "", "Log deletes to this file")
	list := flag.Bool("list", false, "List files only")
	archive := flag.String("archive", "", "Archive directory")
	del := flag.Bool("del", false, "Delete files")
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	mod := flag.String("mod", "", "Modified before")
	flag.Parse()

	var (
		f   = os.Stdout
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext:     *ext,
		size:    *size,
		list:    *list,
		del:     *del,
		wLog:    f,
		archive: *archive,
		mod:     *mod,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)

	if !cfg.restore {
		return filepath.Walk(root,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if filterOut(path, cfg.ext, cfg.size, info, cfg.mod) {
					return nil
				}

				if cfg.list {
					return listFile(path, out)
				}

				if cfg.archive != "" {
					if err := archiveFile(cfg.archive, root, path); err != nil {
						return err
					}
				}

				if cfg.del {
					return delFile(path, delLogger)
				}

				return listFile(path, out)
			})
	}

	return filepath.Walk(cfg.archive,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, "", 0, info, "") {
				return nil
			}

			if err := restore(cfg.archive, root, path); err != nil {
				return err
			}

			return listFile(path, out)
		})
}
