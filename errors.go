package main

import (
	"fmt"
)

// TODO: remove global vars
var hadError = false
var hadRuntimeError = false

func reporter(line int, where string, message string) string {
	s := fmt.Sprintf("[line %v] Error %s: %s\n", line, where, message)
	return s
}

type LexError struct {
	line    int
	message string
}

func NewLexError(line int, message string) error {
	hadError = true
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
	hadError = true
	return &ParseError{token, message}
}

func (err *ParseError) Error() string {
	t := err.token
	if t.typ == EOF {
		return reporter(t.line, " at end", err.message)
	}
	return reporter(t.line, "at '"+t.lexeme+"'", err.message)
}

type RuntimeError struct {
	token   *Token
	message string
}

func NewRuntimeError(token *Token, message string) error {
	hadRuntimeError = true
	return &RuntimeError{token, message}
}

func (err *RuntimeError) Error() string {
	t := err.token
	return reporter(t.line, "at '"+t.lexeme+"'", err.message)
}
