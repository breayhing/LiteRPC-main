package main

import (
	"strings"
	"time"
)

type Stringservice struct{}

type StringArgs struct {
	String1    string
	String2    string
	HandleTime float32
}

type stringMode int

const (
	StringModeCat stringMode = iota
	StringModeCompare
	StringModeTotalLength
	StringModeToUpper
	StringModeToLower
)

var stringArg = &StringArgs{
	String1:    "hELLo",
	String2:    "wORLd",
	HandleTime: 1,
}
var stringRet string

// 添加方法要保证首字母大写
func (f *Stringservice) ToUpper(arg StringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = strings.ToUpper(totalString)
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Stringservice) ToLower(arg StringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = strings.ToLower(totalString)
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Stringservice) Concat(arg StringArgs, reply *string) error {
	*reply = arg.String1 + arg.String2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

func (f *Stringservice) Compare(arg StringArgs, reply *string) error {
	if arg.String1 > arg.String2 {
		*reply = arg.String1
	} else {
		*reply = arg.String2
	}
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}
