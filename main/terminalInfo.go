package main

import (
	"fmt"
	"os"
)

var serverHelp = `
Usage: go run . [options]
Options:
  -h, --help    Show help message
  -l, --listen  Server listen ip
  -p, --port    Server listen port
  -m, --mathMode Server run math mode
	1: add   simply add two numbers
	2: sub	 simply sub two numbers
	3: mul	 simply mul two numbers
	4: div	 simply div two numbers
	5: mod	 simply mod two numbers
  -s, --stringMode Server run string mode
	1: concat   concat two strings
	2: compare  compare two strings and return the bigger one	
	3: length	add two strings and return the length of the result string 
	4: tolower  add two strings and return the result string with all lower case
	5: toupper  add two strings and return the result string with all upper case
  example:
  	go run . -p 8081 -l localhost -m 1
  	it means server listen localhost:8081 and run math mode add, arg is defined in mathMethod.go
`

func terminalMessagePrint() {
	var argc = len(os.Args) - 1
	fmt.Println("命令行参数数量:", argc) //不包括初始执行路径
	for k, v := range os.Args[1:] {
		fmt.Printf("args[%v]=[%v]\n", k, v)
	}
	terminalFunc(argc)
	return
}

// 用于根据参数对应调用不同的方法
func terminalFunc(argc int) {
	if argc == 0 {
		fmt.Println("wrong useage")
		fmt.Fprint(os.Stderr, serverHelp)
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprint(os.Stderr, serverHelp)
	}
}
