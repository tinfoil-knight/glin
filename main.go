package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var hadError = false

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
	tokens := strings.Split(source, " ")
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %v] Error %s: %s\n", line, where, message)
	hadError = true
}
