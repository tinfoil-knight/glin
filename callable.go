package main

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, args []interface{}) interface{}
}
