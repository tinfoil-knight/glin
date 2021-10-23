package main

import "strconv"

type TokenList struct {
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

func New(source string) *TokenList {
	return &TokenList{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

// TODO: check bufio.NewScanner custom Split method

func (tl *TokenList) ScanTokens() *[]Token {
	for !tl.isAtEnd() {
		tl.start = tl.current
		tl.scanToken()
	}
	tl.tokens = append(tl.tokens, Token{EOF, "", nil, tl.line})
	return &tl.tokens
}

func (tl *TokenList) isAtEnd() bool {
	return tl.current >= len(tl.source)
}

// consume next char and return it
func (tl *TokenList) advance() byte {
	tl.current++
	return tl.source[tl.current-1]
}

func (tl *TokenList) addToken(typ TokenType, literal interface{}) {
	text := tl.source[tl.start:tl.current]
	newToken := Token{typ, text, literal, tl.line}
	tl.tokens = append(tl.tokens, newToken)
}

func (tl *TokenList) scanToken() {
	c := tl.advance()

	if v, ok := singleCharLexemes[c]; ok {
		tl.addToken(v, nil)
		return
	}
	if v, ok := multiCharLexemes[c]; ok {
		if tl.match('=') {
			tl.addToken(v[0], nil) // multi-char lexeme
		} else {
			tl.addToken(v[1], nil) // single-char lexeme
		}
		return
	}
	switch c {
	case '/':
		if tl.match('/') {
			// single-line comments
			for tl.peek() != '\n' && !tl.isAtEnd() {
				tl.advance()
			}
		} else if tl.match('*') {
			// block comments
			commentClose := tl.peek() == '*' && tl.peekNext() == '/'
			for !commentClose && !tl.isAtEnd() {
				if tl.peek() == '\n' {
					tl.line++
				}
				tl.advance()
			}
			if tl.isAtEnd() {
				error(tl.line, "unterminated block comment")
				return
			}
			// closing */
			tl.advance()
			tl.advance()
		} else {
			// division
			tl.addToken(SLASH, nil)
		}
	case '\n':
		tl.line++
	case ' ', '\r', '\t':
		// whitespace is ignored
	case '"':
		tl.string()
	default:
		if isDigit(c) {
			tl.number()
		} else if isAlpha(c) {
			tl.identifer()
		} else {
			error(tl.line, "unexpected character")
		}
	}
}

// scan a multi-line string literal
func (tl *TokenList) string() {
	for tl.peek() != '"' && !tl.isAtEnd() {
		if tl.peek() == '\n' {
			tl.line++
		}
		tl.advance()
	}

	if tl.isAtEnd() {
		error(tl.line, "unterminated string")
		return
	}

	// closing "
	tl.advance()

	// remove surrounding ""
	value := tl.source[tl.start+1 : tl.current-1]
	tl.addToken(STRING, value)
}

func (tl *TokenList) number() {
	for isDigit(tl.peek()) {
		tl.advance()
	}

	if tl.peek() == '.' && isDigit(tl.peekNext()) {
		tl.advance()
		for isDigit(tl.peek()) {
			tl.advance()
		}
	}
	value, _ := strconv.ParseFloat(tl.source[tl.start:tl.current], 64)
	tl.addToken(NUMBER, value)
}

// identifier : name supplied for variable, function etc.
func (tl *TokenList) identifer() {
	// isAlphaNumeric also supports '_'
	for isAlphaNumeric(tl.peek()) {
		tl.advance()
	}
	text := tl.source[tl.start:tl.current]
	if typ, ok := keywords[text]; ok {
		// reserved keyword
		tl.addToken(typ, "")
	} else {
		tl.addToken(IDENTIFIER, "")
	}
}

func (tl *TokenList) match(expected byte) bool {
	if tl.isAtEnd() || tl.source[tl.current] != expected {
		return false
	}
	tl.advance()
	return true
}

// one char lookahead
func (tl *TokenList) peek() byte {
	if tl.isAtEnd() {
		return '\000'
	}
	return tl.source[tl.current]
}

// two char lookahead
func (tl *TokenList) peekNext() byte {
	if tl.current+1 >= len(tl.source) {
		return '\000'
	}
	return tl.source[tl.current+1]
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
