package main

// Interpreter implements ExprVisitor
type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	i := Interpreter{}
	return &i
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
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
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
