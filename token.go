package main

import "fmt"

type Token struct {
	typ     TokenType
	lexeme  string      // raw substring in source code
	literal interface{} // fixed value: string, numbers etc.
	line    int         // location info
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %s", t.typ, t.lexeme, t.literal)
}
