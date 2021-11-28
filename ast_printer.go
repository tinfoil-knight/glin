package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// AstPrinter implements ExprVisitor, StmtVisitor
type AstPrinter struct{}

func (a *AstPrinter) Print(statements []Stmt) {
	nodes := a.Create(statements)
	b, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	u, _ := unescapeUnicodeCharacters(b)
	fmt.Println(string(u))
}

func (a *AstPrinter) Create(statements []Stmt) interface{} {
	return Node{
		"_type": "Program",
		"body":  a.resolve(statements),
	}
}

type Node map[string]interface{}

/*
 * ExprVisitor implementation
 */

func (a *AstPrinter) visitBinaryExpr(b *Binary) interface{} {
	return Node{
		"_type":    "BinaryExpression",
		"operator": b.operator.lexeme,
		"left":     a.resolveExpr(b.left),
		"right":    a.resolveExpr(b.right),
	}
}

func (a *AstPrinter) visitGroupingExpr(g *Grouping) interface{} {
	return Node{
		"_type":      "GroupExpression",
		"expression": a.resolveExpr(g.expression),
	}
}

func (a *AstPrinter) visitLiteralExpr(l *Literal) interface{} {
	return Node{
		"_type": getLiteralType(l.value) + "Literal",
		"value": fmt.Sprint(l.value),
	}
}

func (a *AstPrinter) visitUnaryExpr(u *Unary) interface{} {
	return Node{
		"_type":    "UnaryExpression",
		"operator": u.operator.lexeme,
		"right":    a.resolveExpr(u.right),
	}
}

func (a *AstPrinter) visitCallExpr(c *Call) interface{} {
	var args []interface{}
	for _, arg := range c.arguments {
		args = append(args, a.resolveExpr(arg))
	}

	return Node{
		"_type":     "CallExpression",
		"callee":    a.resolveExpr(c.callee),
		"arguments": args,
	}
}

func (a *AstPrinter) visitVariableExpr(v *Variable) interface{} {
	// for superclass
	if v.name == (Token{}) {
		return nil
	}
	return Node{
		"_type": "Identifier",
		"name":  v.name.lexeme,
	}
}

func (a *AstPrinter) visitAssignExpr(e *Assign) interface{} {
	return Node{
		"_type": "AssignmentExpression",
		"left": Node{
			"_type": "Identifier",
			"name":  e.name.lexeme,
		},
		"right": a.resolveExpr(e.value),
	}
}

func (a *AstPrinter) visitLogicalExpr(l *Logical) interface{} {
	return Node{
		"_type":    "LogicalExpression",
		"operator": l.operator.lexeme,
		"left":     a.resolveExpr(l.left),
		"right":    a.resolveExpr(l.right),
	}
}

func (a *AstPrinter) visitGetExpr(g *Get) interface{} {
	return Node{
		"_type":  "GetExpression",
		"object": a.resolveExpr(g.object),
		"property": Node{
			"_type": "Identifier",
			"name":  g.name.lexeme,
		},
	}
}

func (a *AstPrinter) visitSetExpr(s *Set) interface{} {
	return Node{
		"_type": "SetExpression",
		"left": Node{
			"_type":  "MemberExpression",
			"object": a.resolveExpr(s.object),
			"property": Node{
				"_type": "Identifier",
				"name":  s.name.lexeme,
			},
		},
		"right": a.resolveExpr(s.value),
	}
}

func (a *AstPrinter) visitSuperExpr(s *Super) interface{} {
	return Node{
		"_type": "MemberExpression",
		"object": Node{
			"_type": "Super",
		},
		"property": s.method.lexeme,
	}
}

func (a *AstPrinter) visitThisExpr(_ *This) interface{} {
	return Node{
		"_type": "This",
	}
}

/*
 * StmtVisitor implementation
 */

func (a *AstPrinter) visitBlockStmt(stmt *Block) interface{} {
	return Node{
		"_type": "BlockStatement",
		"body":  a.resolve(stmt.statements),
	}
}

func (a *AstPrinter) visitClassStmt(stmt *Class) interface{} {
	var methods []interface{}

	for _, method := range stmt.methods {
		kind := METHOD
		if method.name.lexeme == "init" {
			kind = INITIALIZER
		}
		methods = append(methods, a.resolveFunction(method, kind))
	}

	return Node{
		"_type":      "ClassStatement",
		"id":         stmt.name.lexeme,
		"superclass": a.resolveExpr(&stmt.superclass),
		"body":       methods,
	}
}

func (a *AstPrinter) visitExpressionStmt(stmt *Expression) interface{} {
	return Node{
		"_type":      "ExpressionStatement",
		"expression": a.resolveExpr(stmt.expression),
	}
}

func (a *AstPrinter) visitFunctionStmt(stmt *Function) interface{} {
	return a.resolveFunction(*stmt, FUNCTION)
}

func (a *AstPrinter) visitIfStmt(stmt *If) interface{} {
	return Node{
		"_type":      "IfStatement",
		"condition":  a.resolveExpr(stmt.condition),
		"consequent": a.resolveStmt(stmt.thenBranch),
		"alternate":  a.resolveStmt(stmt.elseBranch),
	}
}

func (a *AstPrinter) visitPrintStmt(stmt *Print) interface{} {
	return Node{
		"_type":      "PrintStatement",
		"expression": a.resolveExpr(stmt.expression),
	}
}

func (a *AstPrinter) visitReturnStmt(stmt *Return) interface{} {
	return Node{
		"_type":    "ReturnStatement",
		"argument": a.resolveExpr(stmt.value),
	}
}

func (a *AstPrinter) visitBreakStmt(_ *Break) interface{} {
	return Node{
		"_type": "BreakStatement",
	}
}

func (a *AstPrinter) visitVarStmt(stmt *Var) interface{} {
	return Node{
		"_type": "VariableDeclaration",
		"id":    stmt.name.lexeme,
		"init":  a.resolveExpr(stmt.initializer),
	}
}

func (a *AstPrinter) visitWhileStmt(stmt *While) interface{} {
	return Node{
		"_type":     "WhileStatement",
		"condition": a.resolveExpr(stmt.condition),
		"body":      a.resolveStmt(stmt.body),
	}
}

/*
 * Utility Methods
 */

func (a *AstPrinter) resolve(statements []Stmt) []interface{} {
	var arr []interface{}
	for _, stmt := range statements {
		x := a.resolveStmt(stmt)
		arr = append(arr, x)
	}
	return arr
}

func (a *AstPrinter) resolveStmt(stmt Stmt) interface{} {
	if stmt != nil {
		return stmt.accept(a).(Node)
	}
	return nil
}

func (a *AstPrinter) resolveExpr(expr Expr) interface{} {
	if expr != nil {
		expr.accept(a)
	}
	return nil
}

func (a *AstPrinter) resolveFunction(f Function, kind FunctionType) interface{} {
	params := []Node{}
	for _, param := range f.params {
		node := Node{
			"_type": "Identifier",
			"name":  param.lexeme,
		}
		params = append(params, node)
	}

	return Node{
		"_type":  "FunctionStatement",
		"id":     f.name.lexeme,
		"kind":   kind,
		"params": params,
		"body":   a.resolve(f.body),
	}
}

// Prevents >, < etc. from being escaped
// Adapted from: https://stackoverflow.com/a/51578927/12531621
func unescapeUnicodeCharacters(raw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// see Parser.primary for possible values
func getLiteralType(value interface{}) string {
	switch value.(type) {
	case float64:
		return "Numeric"
	case string:
		return "String"
	case bool:
		return "Boolean"
	case nil:
		return "Nil"
	default:
		// this case won't occur
		return "Invalid"
	}
}
