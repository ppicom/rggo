package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("p", "", "Project directory")
	branch := flag.String("b", "main", "Branch to push to")
	flag.Parse()

	if err := run(*proj, *branch, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj, branch string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	if branch == "" {
		return fmt.Errorf("git branch is required: %w", ErrValidation)
	}

	data, err := os.ReadFile(".pipeline.yml")
	if err != nil {
		return fmt.Errorf("missing .pipeline.yml file: %w", ErrValidation)
	}

	var p pipelineFile
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return fmt.Errorf("unable to unmarshal .pipeline.yml: %w", ErrValidation)
	}

	pipeline := make([]executer, len(p.Pipeline))

	for i, s := range p.Pipeline {
		switch {
		case s.Exception:
			pipeline[i] = newExceptionStep(s.Name, s.Exe, s.Msg, proj, s.Args)
		case s.Timeout > 0:
			pipeline[i] = newTimeoutStep(s.Name, s.Exe, s.Msg, proj, s.Args, time.Duration(s.Timeout)*time.Second)
		default:
			pipeline[i] = newStep(s.Name, s.Exe, s.Msg, proj, s.Args)
		}
	}

	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	done := make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}

			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errCh <- err
				return
			}
		}
		close(done)
	}()

	for {
		select {
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s: exiting: %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}
