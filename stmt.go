// Code generated by "make codegen"; DO NOT EDIT.
package main

type StmtVisitor interface {
	visitBlockStmt(*Block) interface{}
	visitClassStmt(*Class) interface{}
	visitExpressionStmt(*Expression) interface{}
	visitFunctionStmt(*Function) interface{}
	visitIfStmt(*If) interface{}
	visitWhileStmt(*While) interface{}
	visitPrintStmt(*Print) interface{}
	visitReturnStmt(*Return) interface{}
	visitBreakStmt(*Break) interface{}
	visitVarStmt(*Var) interface{}
}

type Stmt interface {
	accept(StmtVisitor) interface{}
}

type Block struct {
	statements []Stmt
}

func (b *Block) accept(visitor StmtVisitor) interface{} {
	return visitor.visitBlockStmt(b)
}

type Class struct {
	name       Token
	superclass Variable
	methods    []Function
}

func (c *Class) accept(visitor StmtVisitor) interface{} {
	return visitor.visitClassStmt(c)
}

type Expression struct {
	expression Expr
}

func (e *Expression) accept(visitor StmtVisitor) interface{} {
	return visitor.visitExpressionStmt(e)
}

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (f *Function) accept(visitor StmtVisitor) interface{} {
	return visitor.visitFunctionStmt(f)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *If) accept(visitor StmtVisitor) interface{} {
	return visitor.visitIfStmt(i)
}

type While struct {
	condition Expr
	body      Stmt
}

func (w *While) accept(visitor StmtVisitor) interface{} {
	return visitor.visitWhileStmt(w)
}

type Print struct {
	expression Expr
}

func (p *Print) accept(visitor StmtVisitor) interface{} {
	return visitor.visitPrintStmt(p)
}

type Return struct {
	keyword Token
	value   Expr
}

func (r *Return) accept(visitor StmtVisitor) interface{} {
	return visitor.visitReturnStmt(r)
}

type Break struct {
	keyword Token
}

func (b *Break) accept(visitor StmtVisitor) interface{} {
	return visitor.visitBreakStmt(b)
}

type Var struct {
	name        Token
	initializer Expr
}

func (v *Var) accept(visitor StmtVisitor) interface{} {
	return visitor.visitVarStmt(v)
}
