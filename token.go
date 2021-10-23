package main

import "fmt"

type Token struct {
	typ     TokenType
	lexeme  string      // raw substring in source code
	literal interface{} // fixed value: string, numbers etc.
	line    int         // location info
}

// uses generated code to work, ref: tokentype_string
func (t Token) String() string {
	return fmt.Sprintf("Token{%s %s}", t.typ, t.lexeme)
}
