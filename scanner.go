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
	if v, ok := singleCharLexemes[c]; ok {
		tl.addToken(v, "")
	}
	error(tl.line, "Unexpected character.")
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
