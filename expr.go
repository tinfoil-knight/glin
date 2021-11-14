// Code generated by "make codegen"; DO NOT EDIT.
package main

type ExprVisitor interface {
	visitAssignExpr(*Assign) interface{}
	visitBinaryExpr(*Binary) interface{}
	visitCallExpr(*Call) interface{}
	visitGetExpr(*Get) interface{}
	visitGroupingExpr(*Grouping) interface{}
	visitLiteralExpr(*Literal) interface{}
	visitLogicalExpr(*Logical) interface{}
	visitSetExpr(*Set) interface{}
	visitThisExpr(*This) interface{}
	visitUnaryExpr(*Unary) interface{}
	visitVariableExpr(*Variable) interface{}
}

type Expr interface {
	accept(ExprVisitor) interface{}
}

type Assign struct {
	name  Token
	value Expr
}

func (a *Assign) accept(visitor ExprVisitor) interface{} {
	return visitor.visitAssignExpr(a)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b *Binary) accept(visitor ExprVisitor) interface{} {
	return visitor.visitBinaryExpr(b)
}

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (c *Call) accept(visitor ExprVisitor) interface{} {
	return visitor.visitCallExpr(c)
}

type Get struct {
	object Expr
	name   Token
}

func (g *Get) accept(visitor ExprVisitor) interface{} {
	return visitor.visitGetExpr(g)
}

type Grouping struct {
	expression Expr
}

func (g *Grouping) accept(visitor ExprVisitor) interface{} {
	return visitor.visitGroupingExpr(g)
}

type Literal struct {
	value interface{}
}

func (l *Literal) accept(visitor ExprVisitor) interface{} {
	return visitor.visitLiteralExpr(l)
}

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (l *Logical) accept(visitor ExprVisitor) interface{} {
	return visitor.visitLogicalExpr(l)
}

type Set struct {
	object Expr
	name   Token
	value  Expr
}

func (s *Set) accept(visitor ExprVisitor) interface{} {
	return visitor.visitSetExpr(s)
}

type This struct {
	keyword Token
}

func (t *This) accept(visitor ExprVisitor) interface{} {
	return visitor.visitThisExpr(t)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u *Unary) accept(visitor ExprVisitor) interface{} {
	return visitor.visitUnaryExpr(u)
}

type Variable struct {
	name Token
}

func (v *Variable) accept(visitor ExprVisitor) interface{} {
	return visitor.visitVariableExpr(v)
}
