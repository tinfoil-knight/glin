package main

type LoxFunction struct {
	declaration Function
}

func NewLoxFunction(declaration *Function) *LoxFunction {
	return &LoxFunction{*declaration}
}

func (l *LoxFunction) call(interpreter *Interpreter, args []interface{}) interface{} {
	env := NewEnvironment(interpreter.globals)

	for i := 0; i < len(l.declaration.params); i++ {
		env.define(l.declaration.params[i].lexeme, args[i])
	}

	interpreter.executeBlock(l.declaration.body, env)
	return nil
}

func (l *LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
