package main

import (
	"fmt"
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          *Stack
	currentFunction FunctionType
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter, &Stack{}, NONE}
}

/*
 * StmtVisitor implementation
 */

func (r *Resolver) visitBlockStmt(b *Block) interface{} {
	r.beginScope()
	r.resolve(b.statements)
	r.endScope()
	return nil
}

func (r *Resolver) visitClassStmt(c *Class) interface{} {
	r.declare(c.name)
	r.define(c.name)
	return nil
}

func (r *Resolver) visitVarStmt(v *Var) interface{} {
	r.declare(v.name)
	if v.initializer != nil {
		r.resolveExpr(v.initializer)
	}

	r.define(v.name)
	return nil
}

func (r *Resolver) visitFunctionStmt(stmt *Function) interface{} {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) visitExpressionStmt(stmt *Expression) interface{} {
	r.resolveExpr(stmt.expression)
	return nil
}

func (r *Resolver) visitIfStmt(stmt *If) interface{} {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(stmt *Print) interface{} {
	r.resolveExpr(stmt.expression)
	return nil
}

func (r *Resolver) visitReturnStmt(stmt *Return) interface{} {
	if r.currentFunction == NONE {
		fmt.Println(NewParseError(stmt.keyword, "can't return from top-level code"))
	}

	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
	return nil
}

func (r *Resolver) visitWhileStmt(stmt *While) interface{} {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return nil
}

func (r *Resolver) resolveStmt(statement Stmt) {
	statement.accept(r)
}

func (r *Resolver) resolveFunction(function *Function, typ FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typ

	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(function.body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

/*
 * ExprVisitor implementation
 */

func (r *Resolver) visitVariableExpr(v *Variable) interface{} {
	// TODO: check if working
	if !r.scopes.isEmpty() {
		s := r.scopes.peek().get(v.name.lexeme)
		if s != nil && *s == false {
			fmt.Println(NewParseError(v.name, "can't read local variable in its own initializer"))
		}
	}

	r.resolveLocal(v, v.name)
	return nil
}

func (r *Resolver) visitAssignExpr(a *Assign) interface{} {
	r.resolveExpr(a.value)
	r.resolveLocal(a, a.name)
	return nil
}

func (r *Resolver) visitBinaryExpr(b *Binary) interface{} {
	r.resolveExpr(b.left)
	r.resolveExpr(b.right)
	return nil
}

func (r *Resolver) visitCallExpr(c *Call) interface{} {
	r.resolveExpr(c.callee)

	for _, arg := range c.arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) visitGroupingExpr(g *Grouping) interface{} {
	r.resolveExpr(g.expression)
	return nil
}

func (r *Resolver) visitLiteralExpr(_ *Literal) interface{} {
	return nil
}

func (r *Resolver) visitLogicalExpr(l *Logical) interface{} {
	r.resolveExpr(l.left)
	r.resolveExpr(l.right)
	return nil
}

func (r *Resolver) visitUnaryExpr(u *Unary) interface{} {
	r.resolveExpr(u.right)
	return nil
}

func (r *Resolver) resolveExpr(expression Expr) {
	expression.accept(r)
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		if r.scopes.get(i).containsKey(name.lexeme) {
			r.interpreter.resolve(expr, r.scopes.size()-1-i)
			return
		}
	}
}

// Utility Methods

func (r *Resolver) resolve(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) beginScope() {
	r.scopes.push(Scope{})
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) declare(name Token) {
	if r.scopes.isEmpty() {
		return
	}

	scope := r.scopes.peek()

	// prevents duplicate variable declarations in non-global scopes
	if scope.containsKey(name.lexeme) {
		fmt.Println(NewParseError(name, "already a variable with this name in this scope"))
	}

	scope.put(name.lexeme, false)
}

func (r *Resolver) define(name Token) {
	if r.scopes.isEmpty() {
		return
	}

	r.scopes.peek().put(name.lexeme, true)
}

type Scope map[string]bool

func (m Scope) put(key string, value bool) {
	m[key] = value
}

func (m Scope) get(key string) *bool {
	if v, ok := m[key]; ok {
		return &v
	}
	return nil
}

func (m Scope) containsKey(key string) bool {
	_, ok := m[key]
	return ok
}

type Stack []Scope

func (s *Stack) size() int {
	return len(*s)
}

func (s *Stack) isEmpty() bool {
	return s.size() == 0
}

func (s *Stack) push(item Scope) {
	*s = append(*s, item)
}

func (s *Stack) pop() Scope {
	if s.isEmpty() {
		return nil
	}

	i := s.size() - 1
	elem := s.get(i)
	*s = (*s)[:i]
	return elem
}

func (s *Stack) peek() Scope {
	if s.isEmpty() {
		return Scope{}
	}
	return s.get(s.size() - 1)
}

func (s *Stack) get(i int) Scope {
	return (*s)[i]
}
