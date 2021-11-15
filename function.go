package main

// LoxFunction implements LoxCallable
type LoxFunction struct {
	declaration Function
	closure     Environment
	isInit      bool
}

func NewLoxFunction(declaration *Function, closure *Environment, isInit bool) *LoxFunction {
	return &LoxFunction{*declaration, *closure, isInit}
}

func (l *LoxFunction) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	env := NewEnvironment(&l.closure)

	for i := 0; i < len(l.declaration.params); i++ {
		env.define(l.declaration.params[i].lexeme, args[i])
	}
	// handle return statements
	defer func() {
		if err := recover(); err != nil {
			if returnValue, ok := err.(*ReturnT); ok {
				if l.isInit {
					ret = l.closure.getAt(0, "this")
				} else {
					ret = returnValue.value
				}
			} else {
				panic(err)
			}
		}
	}()

	interpreter.executeBlock(l.declaration.body, env)
	if l.isInit {
		return l.closure.getAt(0, "this")
	}
	return nil
}

func (l *LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(&l.closure)
	environment.define("this", instance)
	return NewLoxFunction(&l.declaration, environment, l.isInit)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
