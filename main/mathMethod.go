package main

import "time"

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

//想法：使用数字来对应调用的方法

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
