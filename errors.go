package main

import (
	"fmt"
)

// TODO: remove global var
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

type ParseError struct {
	token   *Token
	message string
}

func NewParseError(token *Token, message string) error {
	return &ParseError{token, message}
}

func (err *ParseError) Error() string {
	t := err.token
	if t.typ == EOF {
		return reporter(t.line, " at end", err.message)
	}
	return reporter(t.line, "at '"+t.lexeme+"'", err.message)
}
