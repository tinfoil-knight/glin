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

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	fmt.Println(tokens)
	parser := NewParser(*tokens)
	expression := parser.Parse()

	if hadError {
		return
	}

	printer := AstPrinter{}
	printer.Print(*expression)

	in := NewInterpreter()
	in.Interpret(*expression)
}
