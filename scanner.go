package main

type TokenList struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
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
	tl.tokens = append(tl.tokens, Token{EOF, "", "", tl.line})
	return &tl.tokens
}

func (tl *TokenList) isAtEnd() bool {
	return tl.current >= len(tl.source)
}

func (tl *TokenList) scanToken() {
	c := tl.advance()
	singleCharLexemes := map[byte]TokenType{
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
	multiCharLexemes := map[byte][]TokenType{
		'!': {BANG_EQUAL, BANG},
		'=': {EQUAL_EQUAL, EQUAL},
		'<': {LESS_EQUAL, LESS},
		'>': {GREATER_EQUAL, GREATER},
	}

	if v, ok := singleCharLexemes[c]; ok {
		tl.addToken(v, "")
		return
	}
	if v, ok := multiCharLexemes[c]; ok {
		if tl.match('=') {
			tl.addToken(v[0], "")
		} else {
			tl.addToken(v[1], "")
		}
		return
	}
	switch c {
	case '/':
		if tl.match('/') {
			for tl.peek() != '\n' && !tl.isAtEnd() {
				tl.advance()
			}
		} else {
			tl.addToken(SLASH, "")
		}
	case '\n':
		tl.line++
	case ' ', '\r', '\t':
		// ignored
	default:
		error(tl.line, "Unexpected character.")
	}
}

func (tl *TokenList) advance() byte {
	tl.current++
	return tl.source[tl.current]
}

func (tl *TokenList) addToken(typ TokenType, literal string) {
	text := tl.source[tl.start:tl.current]
	newToken := Token{typ, text, literal, tl.line}
	tl.tokens = append(tl.tokens, newToken)
}

func (tl *TokenList) match(expected byte) bool {
	if tl.isAtEnd() || tl.source[tl.current] != expected {
		return false
	}
	tl.current++
	return true
}

func (tl *TokenList) peek() byte {
	if tl.isAtEnd() {
		return '\000'
	}
	return tl.source[tl.current]
}
