# glin

A tree-walk interpreter written in Go for the Lox language.

## Get Started

Assuming you have the `Go` language and `make` build utility installed, just run `make build` post cloning this repository to create an executable binary.

The REPL might have bugs currently so it's just better to execute Lox programs through a file.

See the `examples` directory and read through [Lox](https://craftinginterpreters.com/the-lox-language.html) to learn writing Lox programs.

```
./glin examples/hello_world.lox
```

## Note to Self

- Most of the files should've ideally been under a specific sub-package lox but the folder structure is not going to be refactored to preserve the version control history for personal future reference.
- This implementation doesn't have the `clock` native function as added by the author of the book.
- Extensions Implemented:
  - C-style Block Comments (without nesting)
  - REPL automatically prints the results for single expressions
  - `+` operand supports concatenation of string and number
  - break statements

## Attribution

- [Crafting Interpreters](https://craftinginterpreters.com/) by [Robert Nystrom](https://github.com/munificent)
- [AST Explorer](https://astexplorer.net/) was a great resource for exploring ASTs of various languages while I was creating the `AstPrinter` utility.

## References:

- For any non-descript exit code in `os.Exit()`, refer: [UNIX sysexits.h](https://www.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html)
