package main

import (
	"fmt"
	"os"
)

var err error

func methodType() {
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
	for i := 0; i < 3; i++ {
		mathArg.Num1 = i * 3
		mathArg.Num2 = i
		mathArg.HandleTime = 0
		for j := 0; j < Clientnum; i++ {
			//这里客户端调用服务端的方法
			err = Clients[j].Call(Ctx, "Mathservice."+function, mathArg, &mathRet)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("return value:", mathRet)
		}
		fmt.Println("return value:", mathRet)
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
	for i := 0; i < 3; i++ {
		stringArg.HandleTime = 0
		for j := 0; j < Clientnum; i++ {
			//这里客户端调用服务端的方法
			err = Clients[j].Call(Ctx, "Stringservice."+function, stringArg, &stringRet)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("return value:", stringRet)
		}
		fmt.Println("return value:", stringRet)
	}
}
