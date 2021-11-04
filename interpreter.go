package main

import (
	"fmt"
)

// Interpreter implements ExprVisitor, StmtVisitor
type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	i := Interpreter{NewEnvironment(nil)}
	return &i
}

func (i *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if err := recover(); err != nil {
			if iErr, ok := err.(*RuntimeError); ok {
				fmt.Println(iErr)
			} else {
				panic(err)
			}
		}
	}()
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(s Stmt) {
	s.accept(i)
}

func (i *Interpreter) evaluate(e Expr) interface{} {
	return e.accept(i)
}

/*
 * ExprVisitor implementation
 */

func (i *Interpreter) visitLiteralExpr(l *Literal) interface{} {
	return l.value
}

func (i *Interpreter) visitLogicalExpr(l *Logical) interface{} {
	left := i.evaluate(l.left)

	if l.operator.typ == OR {
		if isTruthy(left) {
			return left
		}
	} else {
		// l.operator.typ == AND
		if !isTruthy(left) {
			return left
		}
	}
	return i.evaluate(l.right)
}

func (i *Interpreter) visitGroupingExpr(g *Grouping) interface{} {
	return i.evaluate(g.expression)
}

func (i *Interpreter) visitUnaryExpr(u *Unary) interface{} {
	right := i.evaluate(u.right)

	switch u.operator.typ {
	case MINUS:
		checkNumberOperand(&u.operator, right)
		return -(right).(float64)
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func (i *Interpreter) visitBinaryExpr(b *Binary) interface{} {
	left := i.evaluate(b.left)
	right := i.evaluate(b.right)

	switch b.operator.typ {
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}

		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		panic(NewRuntimeError(&b.operator, "operands must be two numbers or two strings"))
	case MINUS:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(&b.operator, left, right)
		return left.(float64) <= right.(float64)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case BANG_EQUAL:
		return !isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) visitVariableExpr(v *Variable) interface{} {
	return i.env.get(&v.name)
}

func (i *Interpreter) visitAssignExpr(a *Assign) interface{} {
	value := i.evaluate(a.value)
	i.env.assign(&a.name, value)
	return value
}

func isTruthy(v interface{}) bool {
	switch v.(type) {
	case nil:
		return false
	case bool:
		return v.(bool)
	default:
		return true
	}
}

func isEqual(a interface{}, b interface{}) bool {
	return a == b
}

func checkNumberOperand(operator *Token, value interface{}) {
	if _, ok := value.(float64); ok {
		return
	}
	panic(NewRuntimeError(operator, "operand must be a number"))
}

func checkNumberOperands(operator *Token, left interface{}, right interface{}) {
	checkNumberOperand(operator, left)
	checkNumberOperand(operator, right)
}

/*
 * StmtVisitor implementation
 */

func (i *Interpreter) visitExpressionStmt(stmt *Expression) interface{} {
	i.evaluate(stmt.expression)
	return nil
}

func (i *Interpreter) visitPrintStmt(stmt *Print) interface{} {
	v := i.evaluate(stmt.expression)
	fmt.Println(v)
	return nil
}

func (i *Interpreter) visitIfStmt(stmt *If) interface{} {
	if isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else {
		if stmt.elseBranch != nil {
			i.execute(stmt.elseBranch)
		}
	}
	return nil
}

func (i *Interpreter) visitVarStmt(stmt *Var) interface{} {
	// variables are initialized to nil if value is not provided
	var value interface{}

	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}

	i.env.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) visitBlockStmt(stmt *Block) interface{} {
	i.executeBlock(stmt.statements, NewEnvironment(i.env))
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.env

	defer func() { i.env = previous }()

	i.env = environment

	for _, statement := range statements {
		i.execute(statement)
	}
}
