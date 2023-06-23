package main

type stringService struct{}

type stringArgs struct {
	Num1       int
	Num2       int
	HandleTime float32
}

type stringMode int

const (
	StringModeCat stringMode = iota
)
