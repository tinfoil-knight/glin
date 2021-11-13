package main

// LoxClass implements LoxCallable
type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name}
}

func (l *LoxClass) arity() int {
	return 0
}

func (l *LoxClass) call(_ *Interpreter, _ []interface{}) interface{} {
	instance := LoxInstance{*l}
	return instance
}

func (l *LoxClass) String() string {
	return l.name
}

type LoxInstance struct {
	class LoxClass
}

func (l LoxInstance) String() string {
	return l.class.name + " instance"
}
