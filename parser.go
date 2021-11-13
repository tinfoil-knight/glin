package main

import "fmt"

// Parser implements ExprVisitor, StmtVisitor
// Uses recursive descent parsing
type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if err := recover(); err != nil {
			if pErr, ok := err.(ParseError); ok {
				fmt.Println(pErr)
				p.synchronize()
			} else {
				panic(err)
			}
		}
	}()

	switch true {
	case p.match(CLASS):
		return p.classDeclaration()
	case p.match(FUN):
		return p.function("function")
	case p.match(VAR):
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "expect class name")
	p.consume(LEFT_BRACE, "expect '{' before class body")

	methods := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}

	p.consume(RIGHT_BRACE, "expect '}' after class body")

	return &Class{name, methods}
}

// @param kind: "function", "method"
func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, "expect "+kind+" name")

	p.consume(LEFT_PAREN, "expect '(' after "+kind+" name")

	var parameters []Token

	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				fmt.Println(NewParseError(p.peek(), "can't have more than 255 parameters"))
			}

			parameters = append(parameters, p.consume(IDENTIFIER, "expect parameter name"))
			if !p.match(COMMA) {
				break
			}
		}
	}

	p.consume(RIGHT_PAREN, "expect ')' after parameters")

	p.consume(LEFT_BRACE, "expect '{' before "+kind+" body")

	body := p.block()

	return &Function{name, parameters, body}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "expect variable name")
	var initializer Expr

	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "expect ';' after variable declaration")
	return &Var{name, initializer}
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(PRINT):
		return p.printStatement()
	case p.match(RETURN):
		return p.returnStatement()
	case p.match(IF):
		return p.ifStatement()
	case p.match(FOR):
		return p.forStatement()
	case p.match(WHILE):
		return p.whileStatement()
	case p.match(LEFT_BRACE):
		return &Block{p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "expect '}' after block")
	return statements
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "expect ';' after value")
	return &Print{value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()

	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}

	p.consume(SEMICOLON, "expect ';' after return value")
	return &Return{keyword, value}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'if'")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "expect ')' after condition")

	thenBranch := p.statement()
	elseBranch := (Stmt)(nil)

	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &If{condition, thenBranch, elseBranch}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'for'")

	initializer := (Stmt)(nil)
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	condition := (Expr)(nil)
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}

	p.consume(SEMICOLON, "expect ';' after loop condition")

	increment := (Expr)(nil)
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}

	p.consume(RIGHT_PAREN, "expect ')' after for clauses")

	body := p.statement()

	if increment != nil {
		body = &Block{[]Stmt{body, &Expression{increment}}}
	}

	if condition == nil {
		condition = &Literal{true}
	}

	body = &While{condition, body}

	if initializer != nil {
		body = &Block{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'while'")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "expect ')' after condition")

	body := p.statement()

	return &While{condition, body}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "expect ';' after expression")
	return &Expression{expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if e, ok := (expr).(*Variable); ok {
			name := e.name
			return &Assign{name, value}
		} else if e, ok := (expr).(*Get); ok {
			return &Set{e.object, e.name, value}
		}

		fmt.Println(NewParseError(equals, "invalid assignment target"))
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparision()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparision()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparision() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
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

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "expect property name after '.'")
			expr = &Get{expr, name}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return &Literal{false}
	case p.match(TRUE):
		return &Literal{true}
	case p.match(NIL):
		return &Literal{nil}
	case p.match(NUMBER, STRING):
		return &Literal{p.previous().literal}
	case p.match(IDENTIFIER):
		return &Variable{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "expect ')' after expression.")
		return &Grouping{expr}
	}

	panic(NewParseError(p.peek(), "expect expression"))
}

// match checks if current token matches any given type
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// check returns true for token of given type
func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().typ == typ
}

// advance consumes and returns current token
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd returns true if no tokens are left
func (p *Parser) isAtEnd() bool {
	return p.peek().typ == EOF
}

// peek returns current token yet to consume
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// previous returns last consumed token
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// consume looks for given token, panics if not found
func (p *Parser) consume(t TokenType, msg string) Token {
	if p.check(t) {
		return p.advance()
	}
	panic(NewParseError(p.peek(), msg))
}

// synchronize discards token unless at a statement boundary
// restores state on parsing errors
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().typ == SEMICOLON {
			return
		}

		switch p.peek().typ {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT:
			// discard tokens
		case RETURN:
			return
		}

		p.advance()
	}
}

// finishCall returns a Call AST node with 0 or more arguments
func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}

	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				fmt.Println(NewParseError(p.peek(), "can't have more than 255 arguments"))
			}
			arguments = append(arguments, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}

	// location context in errors
	paren := p.consume(RIGHT_PAREN, "expect ')' after arguments")

	return &Call{callee, paren, arguments}
}
