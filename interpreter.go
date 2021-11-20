package main

import (
	"fmt"
)

// Interpreter implements ExprVisitor, StmtVisitor
type Interpreter struct {
	env      *Environment
	globals  *Environment
	locals   map[Expr]int
	replMode bool
}

func NewInterpreter(replMode bool) *Interpreter {
	globals := NewEnvironment(nil)
	env := *globals
	locals := map[Expr]int{}
	i := Interpreter{&env, globals, locals, replMode}
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
		if !i.replMode {
			i.execute(stmt)
		} else {
			if v, ok := (stmt).(*Expression); ok {
				fmt.Println(i.evaluate(v.expression))
			} else {
				i.execute(stmt)
			}
		}
	}
}

func (i *Interpreter) execute(s Stmt) {
	s.accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
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

func (i *Interpreter) visitSetExpr(s *Set) interface{} {
	object := i.evaluate(s.object)
	v, ok := object.(LoxInstance)

	if !ok {
		panic(NewRuntimeError(s.name, "only instances have fields"))
	}

	value := i.evaluate(s.value)
	v.set(s.name, value)
	return value
}

func (i *Interpreter) visitSuperExpr(s *Super) interface{} {
	distance := i.locals[s]
	superclass := i.env.getAt(distance, "super").(*LoxClass)
	object := (i.env.getAt(distance-1, "this")).(*LoxInstance)
	method := superclass.findMethod(s.method.lexeme)
	if method == nil {
		msg := fmt.Sprintf("undefined property %q", s.method.lexeme)
		panic(NewRuntimeError(s.method, msg))
	}
	return method.bind(object)
}

func (i *Interpreter) visitThisExpr(t *This) interface{} {
	v := i.lookUpVariable(t.keyword, t).(*LoxInstance)
	return *v
}

func (i *Interpreter) visitGroupingExpr(g *Grouping) interface{} {
	return i.evaluate(g.expression)
}

func (i *Interpreter) visitUnaryExpr(u *Unary) interface{} {
	right := i.evaluate(u.right)

	switch u.operator.typ {
	case MINUS:
		checkNumberOperand(u.operator, right)
		return -(right).(float64)
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func (i *Interpreter) visitCallExpr(c *Call) interface{} {
	callee := i.evaluate(c.callee)

	var arguments []interface{}

	for _, a := range c.arguments {
		arguments = append(arguments, i.evaluate(a))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		panic(NewRuntimeError(c.paren, "can only call functions and classes"))
	}

	if len(arguments) != function.arity() {
		msg := fmt.Sprintf("expected %d arguments but got %d", function.arity(), len(arguments))
		panic(NewRuntimeError(c.paren, msg))
	}

	return function.call(i, arguments)
}

func (i *Interpreter) visitGetExpr(g *Get) interface{} {
	object := i.evaluate(g.object)

	if v, ok := object.(LoxInstance); ok {
		return v.get(g.name)
	}

	panic(NewRuntimeError(g.name, "only instances have properties"))
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
		panic(NewRuntimeError(b.operator, "operands must be two numbers or two strings"))
	case MINUS:
		checkNumberOperands(b.operator, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		checkNumberOperands(b.operator, left, right)
		// returns +Inf or -Inf on division by zero since all numbers are float64
		return left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(b.operator, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		checkNumberOperands(b.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(b.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(b.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(b.operator, left, right)
		return left.(float64) <= right.(float64)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case BANG_EQUAL:
		return !isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) visitAssignExpr(a *Assign) interface{} {
	value := i.evaluate(a.value)
	if distance, ok := i.locals[a]; ok {
		i.env.assignAt(distance, a.name, value)
	} else {
		i.globals.assign(a.name, value)
	}
	return value
}

func (i *Interpreter) visitVariableExpr(v *Variable) interface{} {
	return i.lookUpVariable(v.name, v)
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.env.getAt(distance, name.lexeme)
	}

	return i.globals.get(name)
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

func checkNumberOperand(operator Token, value interface{}) {
	if _, ok := value.(float64); ok {
		return
	}
	panic(NewRuntimeError(operator, "operand must be a number"))
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) {
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

type ReturnT struct {
	value interface{}
}

func (i *Interpreter) visitReturnStmt(stmt *Return) interface{} {
	value := (interface{})(nil)
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}
	panic(ReturnT{value})
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

func (i *Interpreter) visitWhileStmt(stmt *While) interface{} {
	for isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
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

func (i *Interpreter) visitFunctionStmt(stmt *Function) interface{} {
	function := NewLoxFunction(stmt, i.env, false)
	i.env.define(stmt.name.lexeme, function)
	return nil
}

func (i *Interpreter) visitClassStmt(stmt *Class) interface{} {
	var superclass interface{}

	hasSuperclass := stmt.superclass != (Variable{})

	if hasSuperclass {
		superclass = i.evaluate(&stmt.superclass)
		if _, ok := superclass.(*LoxClass); !ok {
			panic(NewRuntimeError(stmt.superclass.name, "superclass must be a class"))
		}
	}

	i.env.define(stmt.name.lexeme, nil)

	if hasSuperclass {
		i.env = NewEnvironment(i.env)
		i.env.define("super", superclass)
	}

	methods := map[string]LoxFunction{}

	for _, method := range stmt.methods {
		isInit := method.name.lexeme == "init"
		function := NewLoxFunction(&method, i.env, isInit)
		methods[method.name.lexeme] = *function
	}

	if hasSuperclass {
		i.env = i.env.enclosing
	}

	s, _ := superclass.(*LoxClass)
	class := NewLoxClass(stmt.name.lexeme, s, methods)
	i.env.assign(stmt.name, class)
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
