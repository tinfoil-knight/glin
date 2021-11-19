package main

import (
	"fmt"
)

// Resolver implements ExprVisitor, StmtVisitor
type Resolver struct {
	interpreter     *Interpreter
	scopes          *Stack
	currentFunction FunctionType
	currentClass    ClassType
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
	METHOD
	INITIALIZER
)

type ClassType int

const (
	NONE_CLASS ClassType = iota
	CLASS_TYPE
	SUBCLASS_TYPE
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter, &Stack{}, NONE, NONE_CLASS}
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
	enclosingClass := r.currentClass
	r.currentClass = CLASS_TYPE

	r.declare(c.name)
	r.define(c.name)

	if c.superclass != nil {
		r.currentClass = SUBCLASS_TYPE
		v := (c.superclass).(*Variable)
		if c.name.lexeme == v.name.lexeme {
			fmt.Println(NewParseError(v.name, "a class can't inherit from itself"))
		}
		r.resolveExpr(v)
		r.beginScope()
		r.scopes.peek().put("super", true)
	}

	r.beginScope()
	r.scopes.peek().put("this", true)

	for _, method := range c.methods {
		m := method.(*Function)
		declaration := METHOD

		isInit := m.name.lexeme == "init"
		if isInit {
			declaration = INITIALIZER
		}
		r.resolveFunction(m, declaration)
	}

	r.endScope()

	if c.superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
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
		if r.currentFunction == INITIALIZER {
			fmt.Println(NewParseError(stmt.keyword, "can't return a value from an initializer"))
		}
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

func (r *Resolver) visitGetExpr(g *Get) interface{} {
	r.resolveExpr(g.object)
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

func (r *Resolver) visitSetExpr(s *Set) interface{} {
	r.resolveExpr(s.value)
	r.resolveExpr(s.object)
	return nil
}

func (r *Resolver) visitSuperExpr(s *Super) interface{} {
	if r.currentClass == NONE_CLASS {
		fmt.Println(NewParseError(s.keyword, "can't use 'super' outside of a class"))
	} else if r.currentClass == SUBCLASS_TYPE {
		fmt.Println(NewParseError(s.keyword, "can't use 'super' in a class with no superclass"))
	}

	r.resolveLocal(s, s.keyword)
	return nil
}

func (r *Resolver) visitThisExpr(t *This) interface{} {
	if r.currentClass == NONE_CLASS {
		fmt.Println(NewParseError(t.keyword, "can't use 'this' outside of a class"))
		return nil
	}
	r.resolveLocal(t, t.keyword)
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
