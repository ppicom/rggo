package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrSignal     = errors.New("received signal")
)

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("step: %q: %s: cause: %v", s.step, s.msg, s.cause)
}

func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}

	return t.step == s.step
}

func (s *stepErr) Unwrap() error {
	return s.cause
}
