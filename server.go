package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	// 配置文件
	Service *ServiceConfig

	// tun文件
	Tun *os.File

	// 消息
	Messages chan Message

	// 网络连接
	Conns []net.Conn
}

func InitServer(service *ServiceConfig) (server *Server, err error) {
	server = &Server{
		Service: service,
	}
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) Init() {
	s.Messages = make(chan Message)
	s.Conns = make([]net.Conn, 0)
}

func (s *Server) Process(conn net.Conn) {
	var err error
	defer conn.Close()

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = conn.Read(message.Buffer)

		if err != nil {
			log.Printf("read from conn failed: %v\n", err)
			break
		}

		message.Name = "in-net"
		s.Messages <- message
	}
}

func (s *Server) TunServer() {
	var err error
	s.Tun, err = MakeTunDevice(s.Service.Device)

	if err != nil {
		log.Printf("ssssss %v\n", err)
		return
	}

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = s.Tun.Read(message.Buffer)
		if err != nil {
			panic(err)
		}
		message.Name = "in-tun"
		s.Messages <- message
	}
}

func (s *Server) TcpServer() {
	listen, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", s.Service.Listen, s.Service.Port))
	if err != nil {
		logger.Println("helloworld")
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed: err: %v\n", err)
			continue
		}
		s.Conns = append(s.Conns, conn)

		go s.Process(conn)
	}
}

func (s *Server) Run() {
	go s.TcpServer()

	go s.TunServer()

	for {
		select {
		case message := <-s.Messages:
			switch message.Name {
			case "in-net":
				s.Tun.Write(message.Buffer[:message.Len])
			case "in-tun":
				for _, client := range s.Conns {
					client.Write(message.Buffer[:message.Len])
				}
			}
		}
	}
}
