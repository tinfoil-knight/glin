package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var singleCharLexemes = map[byte]TokenType{
	'(': LEFT_PAREN,
	')': RIGHT_PAREN,
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	',': COMMA,
	'.': DOT,
	'-': MINUS,
	'+': PLUS,
	';': SEMICOLON,
	'*': STAR,
}

// lexemes that can have either 1 or 2 chars
var multiCharLexemes = map[byte][]TokenType{
	'!': {BANG_EQUAL, BANG},
	'=': {EQUAL_EQUAL, EQUAL},
	'<': {LESS_EQUAL, LESS},
	'>': {GREATER_EQUAL, GREATER},
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

// TODO: check bufio.NewScanner custom Split method

func (sc *Scanner) ScanTokens() *[]Token {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.scanToken()
	}
	sc.tokens = append(sc.tokens, Token{EOF, "", nil, sc.line})
	return &sc.tokens
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

// consume next char and return it
func (sc *Scanner) advance() byte {
	sc.current++
	return sc.source[sc.current-1]
}

func (sc *Scanner) addToken(typ TokenType, literal interface{}) {
	text := sc.source[sc.start:sc.current]
	newToken := Token{typ, text, literal, sc.line}
	sc.tokens = append(sc.tokens, newToken)
}

func (sc *Scanner) scanToken() {
	c := sc.advance()

	if v, ok := singleCharLexemes[c]; ok {
		sc.addToken(v, nil)
		return
	}
	if v, ok := multiCharLexemes[c]; ok {
		if sc.match('=') {
			sc.addToken(v[0], nil) // multi-char lexeme
		} else {
			sc.addToken(v[1], nil) // single-char lexeme
		}
		return
	}
	switch c {
	case '/':
		if sc.match('/') {
			// single-line comments
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else if sc.match('*') {
			// block comments
			commentClose := sc.peek() == '*' && sc.peekNext() == '/'
			for !commentClose && !sc.isAtEnd() {
				if sc.peek() == '\n' {
					sc.line++
				}
				sc.advance()
			}
			if sc.isAtEnd() {
				fmt.Println(NewLexError(sc.line, "unterminated block comment"))
				return
			}
			// closing */
			sc.advance()
			sc.advance()
		} else {
			// division
			sc.addToken(SLASH, nil)
		}
	case '\n':
		sc.line++
	case ' ', '\r', '\t':
		// whitespace is ignored
	case '"':
		sc.string()
	default:
		if isDigit(c) {
			sc.number()
		} else if isAlpha(c) {
			sc.identifer()
		} else {
			fmt.Println(NewLexError(sc.line, "unexpected character"))
		}
	}
}

// scan a multi-line string literal
func (sc *Scanner) string() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}

	if sc.isAtEnd() {
		fmt.Println(NewLexError(sc.line, "unterminated string"))
		return
	}

	// closing "
	sc.advance()

	// remove surrounding ""
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addToken(STRING, value)
}

func (sc *Scanner) number() {
	for isDigit(sc.peek()) {
		sc.advance()
	}

	if sc.peek() == '.' && isDigit(sc.peekNext()) {
		sc.advance()
		for isDigit(sc.peek()) {
			sc.advance()
		}
	}
	value, _ := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)
	sc.addToken(NUMBER, value)
}

// identifier : name supplied for variable, function etc.
func (sc *Scanner) identifer() {
	// isAlphaNumeric also supports '_'
	for isAlphaNumeric(sc.peek()) {
		sc.advance()
	}
	text := sc.source[sc.start:sc.current]
	if typ, ok := keywords[text]; ok {
		// reserved keyword
		sc.addToken(typ, nil)
	} else {
		sc.addToken(IDENTIFIER, nil)
	}
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() || sc.source[sc.current] != expected {
		return false
	}
	sc.advance()
	return true
}

// one char lookahead
func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return '\000'
	}
	return sc.source[sc.current]
}

// two char lookahead
func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return '\000'
	}
	return sc.source[sc.current+1]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
