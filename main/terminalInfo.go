package main

import (
	"fmt"
	"os"
)

type terminalInfo struct {
	ip         string
	port       string
	nodeNum    string
	clientNum  string
	methodType string
	method     string
}

var Info = terminalInfo{
	ip:         "",
	port:       "",
	nodeNum:    "2",
	clientNum:  "2",
	methodType: "math",
	method:     "1",
}

var defaultOption = `
you have chosen the default option
the server will run on math mode and invoke Sub method
the client will run on random port and random ip
the server will run on random port and random ip


`

var serverHelp = `
Usage: go run . [options]
Options:
  -d, --default  Server run default mode
  -h, --help    Show help message
  -l, --listen  Server listen ip ,means the first server ip,the second server ip will be the first server ip + 1, and so on
  -p, --port    Server listen port
  -n, --nodeNum Server node number
  -c, --client Client node number
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
		go run . -p 8081 -l localhost -m 1 -n 2 -c 3
  	it means server listen localhost:8081 and run math mode add, arg is defined in mathMethod.go
	and there will be 3 clients to invoke 2 server

	rpc closing 

`

func terminalMessagePrint() {
	var argc = len(os.Args) - 1
	fmt.Println("命令行参数数量:", argc) //不包括初始执行路径
	fmt.Println(os.Args[1])
	for k, v := range os.Args[1:] {
		fmt.Printf("args[%v]=[%v]\n", k+1, v)
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
	if os.Args[1] == "-h" {
		fmt.Fprint(os.Stderr, serverHelp)
		os.Exit(0)
	}
	if os.Args[1] == "-d" {
		fmt.Fprint(os.Stderr, defaultOption)
		return
	} else if argc != 10 || os.Args[1] != "-p" || os.Args[3] != "-l" || os.Args[7] != "-n" || os.Args[9] != "-c" {
		fmt.Println("wrong useage, too few or too many arguments")
		fmt.Fprint(os.Stderr, serverHelp)
		os.Exit(0)
	}
	if os.Args[5] != "-m" && os.Args[5] != "-s" {
		fmt.Println("wrong useage, wrong method type")
		fmt.Fprint(os.Stderr, serverHelp)
		os.Exit(0)
	}
	//选择具体调用的方法以及节点个数
	terminalCall()

	Info.port = os.Args[2]
	Info.ip = os.Args[4]
	ADDR = Info.ip + ":" + Info.port
	// fmt.Println("initial ADDR:", ADDR)
	Info.methodType = os.Args[5]
	Info.method = os.Args[7]
	Info.nodeNum = os.Args[8]
	Info.clientNum = os.Args[10]
	return
}

func terminalCall() {

}
