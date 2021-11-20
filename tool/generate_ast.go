package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	path := os.Args[0]
	if len(args) != 1 {
		fmt.Println("Usage:", path, "<output-directory>")
		os.Exit(64)
	}
	outputDir := args[0]

	defineAst(outputDir, "Expr", []string{
		"Assign   : name Token, value Expr",
		"Binary   : left Expr, operator Token, right Expr",
		"Call     : callee Expr, paren Token, arguments []Expr",
		"Get      : object Expr, name Token",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Logical  : left Expr, operator Token, right Expr",
		"Set      : object Expr, name Token, value Expr",
		"Super    : keyword Token, method Token",
		"This     : keyword Token",
		"Unary    : operator Token, right Expr",
		"Variable : name Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Block      : statements []Stmt",
		"Class      : name Token, superclass Variable, methods []Function",
		"Expression : expression Expr",
		"Function   : name Token, params []Token, body []Stmt",
		"If         : condition Expr, thenBranch Stmt, " + "elseBranch Stmt",
		"While		: condition Expr, body Stmt",
		"Print      : expression Expr",
		"Return     : keyword Token, value Expr",
		"Var        : name Token, initializer Expr",
	})
}

func defineAst(outputDir string, baseName string, types []string) {
	filePath, _ := filepath.Abs(outputDir + "/" + strings.ToLower(baseName) + ".go")
	f, _ := os.Create(filePath)
	defer f.Close()

	f.WriteString("// Code generated by \"make codegen\"; DO NOT EDIT.\n")
	f.WriteString("package main\n")

	// create visitor

	f.WriteString(getVisitor(baseName, types))
	f.WriteString("\n")

	// generate base interface
	f.WriteString(getInterface(baseName, []string{fmt.Sprintf("accept(%sVisitor) interface{}", baseName)}))
	f.WriteString("\n")

	// generate sub-types
	for _, r := range types {
		structS := splitAndTrim(r, ":")

		name := structS[0]
		fields := splitAndTrim(structS[1], ",")

		f.WriteString(getStruct(name, fields))
		f.WriteString("\n")
		f.WriteString(getStructMethod(name, "accept", baseName))
	}
}

// items in fields are of the format: "fieldName fieldType"
func getStruct(name string, fields []string) string {
	s := fmt.Sprintln("type", name, "struct {")
	for _, t := range fields {
		s += fmt.Sprintln(t)
	}
	s += "}\n"
	return s
}

// items in methods are of the format: "methodName(param1, param2, ...)"
func getInterface(name string, methods []string) string {
	s := fmt.Sprintln("type", name, "interface {")
	for _, m := range methods {
		s += fmt.Sprintln(m)
	}
	s += "}\n"
	return s
}

func getStructMethod(name string, methodName string, baseName string) string {
	methodSig := fmt.Sprintf("%s(visitor %sVisitor)", methodName, baseName)
	arg := strings.ToLower(name)[0]
	s := fmt.Sprintf("func (%c *%s) %s interface{} {\n", arg, name, methodSig)
	s += fmt.Sprintf("return visitor.visit%s%s(%c)\n", name, baseName, arg)
	s += "}\n"
	return s
}

func getVisitor(name string, types []string) string {
	methods := make([]string, len(types))
	for i, t := range types {
		a := splitAndTrim(t, ":")
		methods[i] = fmt.Sprintf("visit%s%s(*%s) interface{}", a[0], name, a[0])
	}
	return getInterface(name+"Visitor", methods)
}

func splitAndTrim(s string, sep string) []string {
	a := strings.Split(s, sep)
	for i, x := range a {
		a[i] = strings.TrimSpace(x)
	}
	return a
}
