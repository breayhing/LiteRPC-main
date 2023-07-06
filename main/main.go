package main

import (
	"SRPC"
	"context"
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

// 使用数组
var Clients [10]*SRPC.Xclient

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
	// log.Println("old port:", Info.port)
	Info.port = strconv.Itoa(newPort)
	// fmt.Println("new port:", Info.port)
	ADDR = Info.ip + ":" + Info.port
	lock.Unlock()

	l, err := net.Listen("tcp", ADDR)
	// log.Println("server address:", ADDR)

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
		Clients[i] = SRPC.NewXClient(SRPC.RoundRobinSelect, addrReg, codeWay)
		time.Sleep(time.Second * 1)
		// Clients[i].Say()
	}

	Ctx, Cancel = context.WithTimeout(context.Background(), time.Second*3)

	methodCall()
	// for i := 0; i < Clientnum; i++ {
	// 	// Clients[i].Say()
	// 	for j := 1; j < 5; j++ {
	// 		mathARG := &MathArgs{Num1: 80, Num2: 20, HandleTime: 0}
	// 		mathRET := mathRet
	// 		err = Clients[i].Call(Ctx, "Mathservice.Sub", mathARG, &mathRET)
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 			continue
	// 		}
	// 		fmt.Println("return value:", mathRET)
	// 	}
	// }

	Cancel()
}
