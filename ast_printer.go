package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a *AstPrinter) Print(e Expr) string {
	return fmt.Sprint(e.accept(a))
}

func (a *AstPrinter) visitBinaryExpr(b *Binary) interface{} {
	return a.parenthesize(b.operator.lexeme, b.left, b.right)
}

func (a *AstPrinter) visitGroupingExpr(g *Grouping) interface{} {
	return a.parenthesize("group", g.expression)
}

func (a *AstPrinter) visitLiteralExpr(l *Literal) interface{} {
	if l.value == nil {
		return "nil"
	}
	return fmt.Sprint(l.value)
}

func (a *AstPrinter) visitUnaryExpr(u *Unary) interface{} {
	return a.parenthesize(u.operator.lexeme, u.right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) interface{} {
	var b strings.Builder
	fmt.Fprintf(&b, "(%s", name)
	for _, expr := range exprs {
		s := expr.accept(a)
		fmt.Fprintf(&b, " %v", s)
	}
	fmt.Fprint(&b, ")")
	return b.String()
}
