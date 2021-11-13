package main

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name}
}

func (l LoxClass) String() string {
	return l.name
}
