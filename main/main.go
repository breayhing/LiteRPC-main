package main

import (
	"SRPC"
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
	server := SRPC.NewServer(time.Second)
	// 注册服务
	_ = server.Register(&Stringservice{})
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
	_ = SRPC.NewRegistry()
	log.Fatal(http.Serve(l, nil))
}

func main() {
	// 用于实现命令行参数的解析
	terminalMessagePrint()

	codeWay := "json" // 这里选择编码方式为二进制字节流
	var err error
	addr0 := make(chan string)
	addr1 := make(chan string)
	addr2 := make(chan string)

	go startRegistry(addr0)
	<-addr0
	addrReg := "http://localhost:9999/SRPC" // 注册中心地址
	go startServer(addr1, addrReg)
	go startServer(addr2, addrReg)
	// 使用持续连接
	time.Sleep(time.Second * 2)                                     //等待服务端启动
	cli := SRPC.NewXClient(SRPC.RoundRobinSelect, addrReg, codeWay) //这里选择使用的负载均衡算法
	time.Sleep(time.Second * 2)                                     // 等待服务端注册完成

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	//这里开始执行rpc调用
	for i := 0; i < 5; i++ {
		mathArg.Num1 = i * 3
		mathArg.Num2 = i
		mathArg.HandleTime = 0
		err = cli.Call(ctx, "Mathservice.Sub", mathArg, &mathRet)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("return value:", mathRet)
	}
	// fmt.Println("second call start")
	// for i := 0; i < 5; i++ {
	// 	stringArg.HandleTime = 0
	// 	err = cli.Call(ctx, "Stringservice.Compare", stringArg, &stringRet)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		continue
	// 	}
	// 	fmt.Println("return value:", stringRet)
	// }

	cancel()
	time.Sleep(time.Second * 2)
}
