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
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	run(string(data))
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}
		run(string(line))
		hadError = false
	}
}

// TODO: put in a session instead of this
var in = NewInterpreter()

func run(source string) {
	tokens := NewScanner(source).ScanTokens()
	statements := NewParser(*tokens).Parse()

	if hadError || hadRuntimeError {
		return
	}

	in.Interpret(statements)
}
