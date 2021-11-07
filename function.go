package main

type LoxFunction struct {
	declaration Function
}

func NewLoxFunction(declaration *Function) *LoxFunction {
	return &LoxFunction{*declaration}
}

func (l *LoxFunction) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	env := NewEnvironment(interpreter.globals)

	for i := 0; i < len(l.declaration.params); i++ {
		env.define(l.declaration.params[i].lexeme, args[i])
	}
	// handle return statements
	defer func() {
		if err := recover(); err != nil {
			if returnValue, ok := err.(*ReturnT); ok {
				ret = returnValue.value
			} else {
				panic(err)
			}
		}
	}()
	interpreter.executeBlock(l.declaration.body, env)
	return nil
}

func (l *LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
