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

type Foo struct{}

type Arg struct {
	Num1       int
	Num2       int
	HandleTime float32
}

func (f *Foo) Double(arg Arg, reply *int) error {
	*reply = arg.Num1 + arg.Num2
	time.Sleep(time.Second * time.Duration(arg.HandleTime))
	return nil
}

// 开启服务端
func startServer(addr chan<- string, addrReg string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("network error", err)
	}
	log.Println("server runs on", l.Addr().String())
	server := LiteRPC.NewServer(time.Second)
	_ = server.Register(&Foo{})
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
	var err error
	addr0 := make(chan string)
	addr1 := make(chan string)
	addr2 := make(chan string)

	go startRegistry(addr0)
	<-addr0
	addrReg := "http://localhost:9999/LiteRPC"
	go startServer(addr1, addrReg)
	go startServer(addr2, addrReg)
	// 使用随机负载均衡
	// 使用持续连接
	cli := LiteRPC.NewXClient(LiteRPC.ConsistentHash, addrReg)
	var ret int
	arg := &Arg{
		Num1: 10,
		Num2: 20,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	fmt.Println("first call start")
	for i := 0; i < 5; i++ {
		arg.Num1 = i
		arg.Num2 = i * 2
		arg.HandleTime = 0.5
		err = cli.Call(ctx, "Foo.Double", arg, &ret)
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
		err = cli.Call(ctx, "Foo.Double", arg, &ret)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("return value:", ret)
	}
	cancel()
	time.Sleep(time.Second * 2)
}