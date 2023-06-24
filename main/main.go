package main

import (
	"LiteRPC"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// 开启服务端
func startServer(addr chan<- string, addrReg string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("network error", err)
	}
	log.Println("server runs on", l.Addr().String())
	server := LiteRPC.NewServer(time.Second)
	_ = server.Register(&Mathservice{})
	_ = server.PostRegistry(addrReg, l)
	server.Accept(l)
}

func startRegistry(addr chan<- string) {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Println("network error", err)
	}
	log.Println("registry runs on", l.Addr().String())
	addr <- l.Addr().String()
	_ = LiteRPC.NewRegistry()
	log.Fatal(http.Serve(l, nil))
}

func main() {
	// 用于实现命令行参数的解析

	terminalMessagePrint()
	codeWay := "json"
	var err error
	addr0 := make(chan string)
	addr1 := make(chan string)
	addr2 := make(chan string)

	go startRegistry(addr0)
	<-addr0
	addrReg := "http://localhost:9999/LiteRPC" // 注册中心地址
	go startServer(addr1, addrReg)
	go startServer(addr2, addrReg)
	// 使用持续连接
	time.Sleep(time.Second * 2)                                           //等待服务端启动
	cli := LiteRPC.NewXClient(LiteRPC.RoundRobinSelect, addrReg, codeWay) //这里选择使用的负载均衡算法
	time.Sleep(time.Second * 2)                                           // 等待服务端注册完成
	var ret int

	arg := &MathArgs{
		Num1: 10,
		Num2: 20,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	fmt.Println("first call start")
	for i := 0; i < 5; i++ {
		arg.Num1 = i
		arg.Num2 = i * 2
		arg.HandleTime = 0.5
		err = cli.Call(ctx, "Mathservice.Double", arg, &ret)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("return value:", ret)
	}
	fmt.Println("second call start")
	for i := 0; i < 5; i++ {
		arg.Num1 = i
		arg.Num2 = i * 2
		arg.HandleTime = 0
		err = cli.Call(ctx, "Mathservice.Double", arg, &ret)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("return value:", ret)
	}
	cancel()
	time.Sleep(time.Second * 2)
}
