package SRPC

import (
	"SRPC/codec"
	"SRPC/consistenthash"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SelectMode int

type xClient struct {
	addrRegistry string
	timeout      time.Duration
	mode         SelectMode
	mu           sync.Mutex
	isClose      bool
	index        int // index for round robin
	r            *rand.Rand
	addrs        []string // available servers (get from registry)
	clients      map[string]*client
	ch           *consistenthash.ConsistentHash
}

const (
	RandomSelect SelectMode = iota
	RoundRobinSelect
	ConsistentHash
)

func NewXClient(s SelectMode, regAddr string, codecWay string) *xClient {
	c := &xClient{
		addrRegistry: regAddr,
		timeout:      60 * time.Second,
		r:            rand.New(rand.NewSource(time.Now().UnixNano())),
		mode:         s,
		isClose:      false,
		clients:      make(map[string]*client),
	}
	c.index = c.r.Intn(math.MaxInt32 - 1)
	c.ch = consistenthash.NewConsistentHash(10)
	go c.getServers(codecWay)
	return c
}

func (xc *xClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for k, cli := range xc.clients {
		_ = cli.Close()
		delete(xc.clients, k)
	}
	xc.isClose = true
	return nil
}

func (xc *xClient) Dial(addr string, typ codec.Type) (err error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	cli, ok := xc.clients[addr]
	if ok && cli.available {
		return nil
	}
	// make sure old connection is closed
	if ok && !cli.available {
		_ = cli.Close()
	}
	// create new client
	cli = NewClient()
	err = cli.Dial(addr, typ)
	if err != nil {
		log.Println("rpc xclient error:", err)
		return err
	}
	xc.clients[addr] = cli
	return nil
}

func (xc *xClient) DialServers(servers []string, co codec.Type) (n int, err error) {
	for _, s := range servers {
		err = xc.Dial(s, co)
		if err == nil {
			n += 1
			xc.addrs = append(xc.addrs, s)
			xc.ch.Add(s)
		}
	}
	return
}

func (xc *xClient) getServers(codecWay string) {
	for {
		resp, err := http.Get(xc.addrRegistry)
		if err == nil {
			serversString := resp.Header.Get("rpc-servers")
			servers := strings.Split(serversString, ",")
			switch codecWay {
			case "gob":
				//这里选择编码方式为gob
				_, _ = xc.DialServers(servers, codec.GobCodec)
			case "json":
				//这里选择编码方式为json
				_, _ = xc.DialServers(servers, codec.JsonCodec)
			}
		}
		if xc.isClose {
			break
		}
		select {
		case <-time.After(xc.timeout):
			continue
		}
	}

}

func (xc *xClient) Call(ctx context.Context, serviceMethod string, argv, replyv interface{}) error {
	// log.Println("replyv:", replyv)
	log.Println("argv:", argv)
	var idx int
	//获取服务地址
	if len(xc.addrs) == 0 {
		return errors.New("rpc xclient error: not server available")
	}
	var addr string
	switch xc.mode {
	case RandomSelect:
		// log.Println("random select")
		idx = xc.r.Intn(math.MaxInt32-1) % len(xc.addrs)
		addr = xc.addrs[idx]
	case RoundRobinSelect:
		idx = xc.index % len(xc.addrs)
		xc.index += 1
		addr = xc.addrs[idx]
		// log.Println("round robin select")
	case ConsistentHash:
		addr = xc.ch.Get(fmt.Sprintf("%v+%s", argv, serviceMethod))
		// log.Println("consistent hash select")
	}
	cli := xc.clients[addr]
	return cli.Call(ctx, serviceMethod, argv, replyv)
}
