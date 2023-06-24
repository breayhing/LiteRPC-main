package main

import (
	"time"
)

type Mathservice struct{}

type MathArgs struct {
	Num1       int
	Num2       int
	HandleTime float32
}

type MathMode int

const (
	MathModeDouble MathMode = iota
	MathModeSub
	MathModeMul
	MathModeDiv
	MathModeMod
)

// 执行参数
var mathArg = &MathArgs{
	Num1:       10,
	Num2:       20,
	HandleTime: 0,
}
var mathRet int //返回值

func (f *Mathservice) Pow(arg MathArgs, reply *int) error {
	*reply = arg.Num1 ^ arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Mathservice) Double(arg MathArgs, reply *int) error {
	*reply = arg.Num1 + arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Mathservice) Sub(arg MathArgs, reply *int) error {
	*reply = arg.Num1 - arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Mathservice) Mul(arg MathArgs, reply *int) error {
	*reply = arg.Num1 * arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Mathservice) Div(arg MathArgs, reply *int) error {
	*reply = arg.Num1 / arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Mathservice) Mod(arg MathArgs, reply *int) error {
	*reply = arg.Num1 % arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}
