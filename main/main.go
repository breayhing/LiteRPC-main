package main

import (
	"SRPC"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 实现自定义地址
var ADDR string

// 节点数量
var Clientnum int

// 使用锁进行临界区保护
var lock sync.Mutex

// 使用切片增添新的客户端
var Clients []SRPC.Xclient

// ctx用于取消rpc调用
var Ctx context.Context
var Cancel context.CancelFunc

// 开启服务端
func startServer(addrReg string, ADDR string) {
	//如果ADDR为空，则默认为本地随机地址
	if ADDR == "" {
		ADDR = ":0"
	}

	//使用锁进行临界区保护
	lock.Lock()
	newPort, _ := strconv.Atoi(Info.port)
	newPort++
	log.Println("old port:", Info.port)
	Info.port = strconv.Itoa(newPort)
	fmt.Println("new port:", Info.port)
	ADDR = Info.ip + ":" + Info.port
	lock.Unlock()
	l, err := net.Listen("tcp", ADDR)
	log.Println("server address:", ADDR)

	if err != nil {
		log.Println("server network error", err)
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
		log.Println("registry network error", err)
	}
	log.Println("registry runs on", l.Addr().String())
	addr <- l.Addr().String()
	_ = SRPC.NewRegistry()
	log.Fatal(http.Serve(l, nil))
}

func main() {

	// 用于实现命令行参数的解析
	terminalMessagePrint()

	ServerNum, _ := strconv.Atoi(Info.nodeNum)
	Clientnum, _ = strconv.Atoi(Info.clientNum)

	codeWay := "json" // 这里选择编码方式为二进制字节流
	var err error
	addr0 := make(chan string)

	go startRegistry(addr0)
	<-addr0
	//<-addr0读取到了注册中心的地址，并赋值给addrReg
	addrReg := "http://localhost:9999/SRPC" // 注册中心地址

	log.Println("server number:", ServerNum)
	for ServerNum > 0 {
		// 启动协程
		go func() {
			startServer(addrReg, ADDR)
		}()
		ServerNum--
	}

	// 使用持续连接
	time.Sleep(time.Second * 2) //等待服务端启动

	log.Println("client number:", Clientnum)
	for i := 0; i < Clientnum; i++ {
		Clients = append(Clients, *SRPC.NewXClient(SRPC.RoundRobinSelect, addrReg, codeWay))
	}

	time.Sleep(time.Second * 2) // 等待服务端注册完成

	Ctx, Cancel = context.WithTimeout(context.Background(), time.Second)

	// //新建一个客户端
	// cli := SRPC.NewXClient(SRPC.RoundRobinSelect, addrReg, codeWay)
	// //这里开始执行rpc调用
	// for i := 0; i < 5; i++ {
	// 	mathArg.Num1 = i * 2
	// 	mathArg.Num2 = i
	// 	mathArg.HandleTime = 0
	// 	//这里客户端调用服务端的方法
	// 	if i != 0 {
	// 		err = cli.Call(Ctx, "Mathservice.Sub", mathArg, &mathRet)
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 			continue
	// 		}
	// 		fmt.Println("return value:", mathRet)
	// 	}
	// }

	// //新建一个客户端
	// cli := SRPC.NewXClient(SRPC.RoundRobinSelect, addrReg, codeWay)
	// //这里开始执行rpc调用
	// mathArg.Num1 = 21
	// mathArg.Num2 = 1
	// mathArg.HandleTime = 0
	// //这里客户端调用服务端的方法
	// err = cli.Call(Ctx, "Mathservice.Sub", mathArg, &mathRet)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println("return value:", mathRet)

	for _, CLIENT := range Clients {
		for i := 0; i < 5; i++ {
			fmt.Println(i)
			mathARG := &MathArgs{Num1: 10, Num2: 20, HandleTime: 0}
			mathRET := &mathRet
			err = CLIENT.Call(Ctx, "Mathservice.Sub", mathARG, &mathRET)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("return value:", mathRET)
		}
	}

	Cancel()
}
