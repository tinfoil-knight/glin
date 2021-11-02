// Code generated by "make codegen"; DO NOT EDIT.
package main

type StmtVisitor interface {
	visitBlockStmt(*Block) interface{}
	visitExpressionStmt(*Expression) interface{}
	visitPrintStmt(*Print) interface{}
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

type Expression struct {
	expression Expr
}

func (e *Expression) accept(visitor StmtVisitor) interface{} {
	return visitor.visitExpressionStmt(e)
}

type Print struct {
	expression Expr
}

func (p *Print) accept(visitor StmtVisitor) interface{} {
	return visitor.visitPrintStmt(p)
}

type Var struct {
	name        Token
	initializer Expr
}

func (v *Var) accept(visitor StmtVisitor) interface{} {
	return visitor.visitVarStmt(v)
}
