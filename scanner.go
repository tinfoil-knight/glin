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
	switch c {
	case '(':
		tl.addToken(LEFT_PAREN, "")
	case ')':
		tl.addToken(RIGHT_PAREN, "")
	case '{':
		tl.addToken(LEFT_BRACE, "")
	case '}':
		tl.addToken(RIGHT_BRACE, "")
	case ',':
		tl.addToken(COMMA, "")
	case '.':
		tl.addToken(DOT, "")
	case '-':
		tl.addToken(MINUS, "")
	case '+':
		tl.addToken(PLUS, "")
	case ';':
		tl.addToken(SEMICOLON, "")
	case '*':
		tl.addToken(STAR, "")
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
