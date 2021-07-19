package main

import "fmt"

type Token struct {
	typ     TokenType
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %s", t.typ, t.lexeme, t.literal)
}
