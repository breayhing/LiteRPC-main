package main

import (
	"fmt"
	"os"
)

var serverHelp = `
Usage: go run main.go [options]

Options:
  -h, --help    Show help message
  -l, --listen  Server listen ip
  -p, --port    Server listen port
  -d, --default Server listen default port:(ip:localhost  ,  port:8081), mathMode add
  -m, --mathMode Server run math mode
	1: add
	2: sub
	3: mul
	4: div
	5: mod
  -s, --stringMode Server run string mode
	1: concat
	2: compare
	3: length
	4: substring
	5: replace
  example:
  	go run . -p 8081 -l localhost -m 1
  	it means server listen localhost:8081 and run math mode add
`

var clientHelp = `
Usage: go run main.go [options]

Options:
  -h, --help    Show help message
  -i, --ip that client connect to
  -p, --port that client connect to
`

func terminalMessagePrint(argc int) {
	fmt.Println("命令行参数数量:", argc) //不包括初始执行路径
	for k, v := range os.Args[1:] {
		fmt.Printf("args[%v]=[%v]\n", k, v)
	}
	if argc == 0 {
		fmt.Println("wrong useage")
		fmt.Fprint(os.Stderr, serverHelp)
	}
	if argc == 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Fprint(os.Stderr, serverHelp)
		}
	}
	return
}

// 用于根据参数对应调用不同的方法
func terminalFunc(argc int) {
}
