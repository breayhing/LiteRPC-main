package main

import "strings"

type Stringservice struct{}

type stringArgs struct {
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

var stringArg = &stringArgs{
	String1:    "hELLo",
	String2:    "wORLd",
	HandleTime: 1,
}
var stringRet string

func (f *Stringservice) Concat(arg stringArgs, reply *string) error {
	*reply = arg.String1 + arg.String2
	return nil
}

func (f *Stringservice) Compare(arg stringArgs, reply *string) error {
	if arg.String1 > arg.String2 {
		*reply = arg.String1
	} else {
		*reply = arg.String2
	}
	return nil
}

func (f *Stringservice) totalLength(arg stringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = string(len(totalString))
	return nil
}

func (f *Stringservice) toUpper(arg stringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = strings.ToUpper(totalString)
	return nil
}

func (f *Stringservice) toLower(arg stringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = strings.ToLower(totalString)
	return nil
}
