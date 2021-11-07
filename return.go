package main

type ReturnT struct {
	value interface{}
}

func NewReturn(value interface{}) *ReturnT {
	return &ReturnT{value}
}
