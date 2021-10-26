package main

import "fmt"

// Parser uses recursive descent parsing
type Parser struct {
	tokens  []Token
	current int // zero value for numeric types is 0
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (ex *Expr) {
	defer func() {
		if err := recover(); err != nil {
			if pErr, ok := err.(*ParseError); ok {
				fmt.Println(pErr)
				ex = nil
			} else {
				panic(err)
			}
		}
	}()
	res := p.expression()
	return &res
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparision()

	if p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparision()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparision() Expr {
	expr := p.term()

	if p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	if p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	if p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &Unary{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &Literal{false}
	}
	if p.match(TRUE) {
		return &Literal{true}
	}
	if p.match(NIL) {
		return &Literal{nil}
	}
	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "expect ')' after expression.")
		return &Grouping{expr}
	}

	a := p.peek()
	panic(NewParseError(&a, "expect expression"))
}

// check if current token matches any given type
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// returns true if token of given type
func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().typ == typ
}

// consumes and returns current token
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// check if no tokens are left
func (p *Parser) isAtEnd() bool {
	return p.peek().typ == EOF
}

// returns current token yet to consume
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// returns last consumed token
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// looks for given token, panics if not found
func (p *Parser) consume(t TokenType, msg string) Token {
	if p.check(t) {
		return p.advance()
	}
	a := p.peek()
	panic(NewParseError(&a, msg))
}

// discards token unless at a statement boundary
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().typ == SEMICOLON {
			return
		}

		switch p.peek().typ {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT:
			// discard token
		case RETURN:
			return
		}

		p.advance()
	}
}
