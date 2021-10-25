package main

import (
	"fmt"
)

var hadError = false

func reporter(line int, where string, message string) string {
	hadError = true
	s := fmt.Sprintf("[line %v] Error %s: %s\n", line, where, message)
	return s
}

type LexError struct {
	line    int
	message string
}

func NewLexError(line int, message string) error {
	return &LexError{line, message}
}

func (err *LexError) Error() string {
	return reporter(err.line, "", err.message)
}
