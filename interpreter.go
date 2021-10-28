package main

import "fmt"

// Interpreter implements ExprVisitor
type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	i := Interpreter{}
	return &i
}

func (i *Interpreter) Interpret(e Expr) {
	defer func() {
		if err := recover(); err != nil {
			if iErr, ok := err.(*RuntimeError); ok {
				fmt.Println(iErr)
			} else {
				panic(err)
			}
		}
	}()
	result := i.evaluate(e)
	fmt.Println(result)
}

func (i *Interpreter) evaluate(e Expr) interface{} {
	return e.accept(i)
}

func (i *Interpreter) visitLiteralExpr(l *Literal) interface{} {
	return l.value
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