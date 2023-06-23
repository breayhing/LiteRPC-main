package main

import "strings"

type stringService struct{}

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

func (f *stringService) Concat(arg stringArgs, reply *string) error {
	*reply = arg.String1 + arg.String2
	return nil
}

func (f *stringService) Compare(arg stringArgs, reply *string) error {
	if arg.String1 > arg.String2 {
		*reply = arg.String1
	} else {
		*reply = arg.String2
	}
	return nil
}

func (f *stringService) totalLength(arg stringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = string(len(totalString))
	return nil
}

func (f *stringService) toUpper(arg stringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = strings.ToUpper(totalString)
	return nil
}

func (f *stringService) toLower(arg stringArgs, reply *string) error {
	totalString := arg.String1 + arg.String2
	*reply = strings.ToLower(totalString)
	return nil
}
