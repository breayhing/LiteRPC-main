package main

import (
	"fmt"
	"os"
)

var err error

func methodCall() {
	var function = os.Args[6]
	if Info.methodType == "-m" {
		fmt.Println("math mode")
		mathCall(function)
	}
	if Info.methodType == "-s" {
		fmt.Println("string mode")
		stringCall(function)
	}

}

func mathCall(function string) {
	switch function {
	case "1":
		function = "Add"
	case "2":
		function = "Sub"
	case "3":
		function = "Mul"
	case "4":
		function = "Div"
	case "5":
		function = "Mod"
	case "6":
		function = "Pow"
	}

	//默认是一个循环进行计算
	for i := 0; i < Clientnum; i++ {
		// Clients[i].Say()
		for j := 1; j < 5; j++ {
			mathARG := &MathArgs{Num1: 80, Num2: 20, HandleTime: 0}
			mathRET := mathRet
			err = Clients[i].Call(Ctx, "Mathservice."+function, mathARG, &mathRET)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("return value:", mathRET)
		}
	}
}

func stringCall(function string) {
	switch function {
	case "1":
		function = "ToUpper"
	case "2":
		function = "ToLower"
	case "3":
		function = "Concat"
	case "4":
		function = "Compare"
	}
	for i := 0; i < Clientnum; i++ {
		// Clients[i].Say()
		for j := 1; j < 5; j++ {
			stringArg := &StringArgs{String1: "hello", String2: "world", HandleTime: 0}
			stringRet := stringRet
			err = Clients[i].Call(Ctx, "Stringservice."+function, stringArg, &stringRet)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("return value:", stringRet)
		}
	}
}
