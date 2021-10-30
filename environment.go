package main

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.put(name, value)
}

func (e *Environment) get(name *Token) interface{} {
	if elem, ok := e.values[name.lexeme]; ok {
		return elem
	}
	panic(NewRuntimeError(name, "undefined variable '"+name.lexeme+"'."))
}

func (e *Environment) put(name string, value interface{}) {
	e.values[name] = value
}
