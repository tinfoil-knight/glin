package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		// Start interactive session
	case 1:
		// Read file
	default:
		fmt.Println("Usage: glin <file-name>")
		os.Exit(64)
	}
}
