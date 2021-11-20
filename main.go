package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args[1:]
	path := os.Args[0]
	switch len(args) {
	case 0:
		runPrompt()
	case 1:
		runFile(args[0])
	default:
		fmt.Println("Usage:", path, "<file-name>")
		os.Exit(64)
	}
}

func runFile(path string) {
	s := NewSession(false)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	run(string(data), s)
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	s := NewSession(true)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}
		run(string(line), s)
		hadError = false
		hadRuntimeError = false
	}
}

func run(source string, s Session) {
	tokens := NewScanner(source).ScanTokens()
	statements := NewParser(tokens).Parse()

	if hadError {
		return
	}

	s.resolver.resolve(statements)

	if hadError || hadRuntimeError {
		return
	}

	s.interpreter.Interpret(statements)
}

type Session struct {
	interpreter *Interpreter
	resolver    *Resolver
}

func NewSession(replMode bool) Session {
	in := NewInterpreter(replMode)

	return Session{
		interpreter: in,
		resolver:    NewResolver(in),
	}
}
