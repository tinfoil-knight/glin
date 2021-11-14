package main

// LoxClass implements LoxCallable
type LoxClass struct {
	name    string
	methods map[string]LoxFunction
}

func NewLoxClass(name string, methods map[string]LoxFunction) *LoxClass {
	return &LoxClass{name, methods}
}

func (l *LoxClass) arity() int {
	return 0
}

func (l *LoxClass) call(_ *Interpreter, _ []interface{}) interface{} {
	instance := LoxInstance{*l, map[string]interface{}{}}
	return instance
}

func (l *LoxClass) findMethod(name string) *LoxFunction {
	if v, ok := l.methods[name]; ok {
		return &v
	}
	return nil
}

func (l LoxClass) String() string {
	return l.name
}

type LoxInstance struct {
	class  LoxClass
	fields map[string]interface{}
}

func (l *LoxInstance) get(name Token) interface{} {
	if v, ok := l.fields[name.lexeme]; ok {
		return v
	}

	method := l.class.findMethod(name.lexeme)
	if method != nil {
		return method.bind(l)
	}

	panic(NewRuntimeError(name, "undefined property '"+name.lexeme+"'."))
}

func (l *LoxInstance) set(name Token, value interface{}) {
	l.fields[name.lexeme] = value
}

func (l LoxInstance) String() string {
	return l.class.name + " instance"
}
