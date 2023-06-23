package LiteRPC

import (
	"LiteRPC/codec"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	handleTimeout time.Duration
	serviceMap    sync.Map
}

func NewServer(timeout time.Duration) *Server {
	return &Server{
		handleTimeout: timeout,
	}
}

func (s *Server) Register(recv interface{}) error {
	serve := NewService(recv)
	_, duplicate := s.serviceMap.LoadOrStore(serve.name, serve)
	if duplicate {
		log.Println("rpc server: service already loaded: " + serve.name)
		return errors.New("rpc server: service already loaded: " + serve.name)
	}
	return nil
}

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error", err)
			return
		}
		go s.ServeConn(conn)
	}
}

func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	//读取option
	opt := make([]byte, 5)
	for readBytes := 0; readBytes < 5; {
		n, err := conn.Read(opt[readBytes:])
		if err != nil {
			if err != io.EOF {
				log.Println("rpc server: conn read error", err)
			}
			break
		}
		readBytes += n
	}
	//进行协议解析
	newCodecFunc, err := codec.ParseOption(opt)
	log.Println("option: ", opt)
	if err != nil {
		log.Println("rpc server: parsing option error")
		return
	}
	log.Println("newCodecFunc: ", newCodecFunc)
	s.ServeCodec(newCodecFunc(conn))
}

func (s *Server) ServeCodec(c codec.Codec) {
	header := new(codec.Header)
	var err error
	for {
		done := make(chan struct{})
		// 读取头部
		err = c.ReadHeader(header)
		if err != nil {
			log.Println("rpc server: parsing header error")
			return
		}
		// 获取服务方法
		serviceMethod := header.ServiceMethod
		serviceMethodStrings := strings.Split(serviceMethod, ".")
		if len(serviceMethodStrings) != 2 {
			log.Println("rpc server: ill formed service method")
			err = errors.New("rpc server: ill formed service method")
			return
		}
		servei, ok := s.serviceMap.Load(serviceMethodStrings[0])
		if !ok {
			log.Println("rpc server: request service unavailable")
			err = errors.New("rpc server: ill formed service method")
			return
		}
		serve := servei.(*service)
		methodTyp := serve.getMethod(serviceMethodStrings[1])
		if methodTyp == nil {
			log.Println("rpc server: request method unavailable")
			err = errors.New("rpc server: request method unavailable")
			return
		}
		go func() {
			defer close(done)
			// 获取输入输出参数
			argv := methodTyp.newArgv()
			replyv := methodTyp.newReplyv()
			body := argv.Addr().Interface()
			err = c.ReadBody(body)
			if err != nil {
				log.Println("rpc server: parsing body error")
				return
			}
			err = serve.call(methodTyp, argv, replyv)
			var replyvi interface{}
			if err != nil {
				log.Println("rpc server: calling error " + err.Error())
				header.Error = err.Error()
				replyvi = nil
			} else {
				replyvi = replyv.Interface()
			}
			// 发送响应
			err = c.Write(header, replyvi)
			if err != nil {
				log.Println("rpc server: write response error")
				return
			}
		}()
		select {
		// 超时处理
		case <-time.After(s.handleTimeout):
			log.Println("rpc server: handle timeout")
			err = errors.New("rpc server: handle timeout")
			continue
		case <-done:
		}
		if err != nil {
			break
		}
	}
}

// 服务注册
func (s *Server) PostRegistry(addrRegistry string, lis net.Listener) error {
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, addrRegistry, nil)
	if err != nil {
		return err
	}
	req.Header.Set("rpc-server-addr", lis.Addr().String())
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}
