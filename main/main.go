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

// 使用锁进行临界区保护
var lock sync.Mutex

var serverseq = false

// 开启服务端
func startServer(addrReg string, ADDR string) {
	//如果ADDR为空，则默认为本地随机地址
	if ADDR == "" {
		ADDR = ":0"
	}
	if serverseq != false {
		lock.Lock()
		newPort, _ := strconv.Atoi(Info.port)
		newPort++
		Info.port = strconv.Itoa(newPort)
		fmt.Println("new port:", Info.port)
		ADDR = Info.ip + ":" + Info.port
		lock.Unlock()
	}
	l, err := net.Listen("tcp", ADDR)
	fmt.Println("listen seems ok")
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

	codeWay := "json" // 这里选择编码方式为二进制字节流
	var err error
	addr0 := make(chan string)

	go startRegistry(addr0)
	<-addr0
	//<-addr0读取到了注册中心的地址，并赋值给addrReg
	addrReg := "http://localhost:9999/SRPC" // 注册中心地址
	go startServer(addrReg, ADDR)
	serverseq = true
	go startServer(addrReg, ADDR)
	// 使用持续连接
	time.Sleep(time.Second * 2) //等待服务端启动

	//这里选择使用的负载均衡算法，此处为轮询算法

	//新建一个客户端
	cli := SRPC.NewXClient(SRPC.RoundRobinSelect, addrReg, codeWay)
	client2 := SRPC.NewXClient(SRPC.ConsistentHash, addrReg, codeWay)

	time.Sleep(time.Second * 2) // 等待服务端注册完成

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	//这里开始执行rpc调用
	for i := 0; i < 5; i++ {
		mathArg.Num1 = i * 3
		mathArg.Num2 = i
		mathArg.HandleTime = 0
		//这里客户端调用服务端的方法
		err = cli.Call(ctx, "Mathservice.Sub", mathArg, &mathRet)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("return value:", mathRet)
	}
	fmt.Println("second call start")
	for i := 0; i < 5; i++ {
		stringArg.HandleTime = 0
		err = client2.Call(ctx, "Stringservice.Compare", stringArg, &stringRet)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("return value:", stringRet)
	}

	cancel()
	time.Sleep(time.Second * 2)
}
