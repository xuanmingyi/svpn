package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Client struct {
	// 配置文件
	Service *ServiceConfig

	// tun文件
	Tun *os.File

	// 消息
	Messages chan Message

	// 网络连接
	Conn net.Conn
}

func InitClient(service *ServiceConfig) (client *Client, err error) {
	client = &Client{
		Service: service,
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Init() {
	c.Messages = make(chan Message)
}

func (c *Client) TunServer() {
	var err error

	c.Tun, err = MakeTunDevice(c.Service.Device)

	if err != nil {
		log.Printf("ssss: %v\n", err)
		return
	}

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = c.Tun.Read(message.Buffer)
		if err != nil {
			panic(err)
		}
		message.Name = "in-tun"
		c.Messages <- message
	}
}

func (c *Client) TcpClient() {
	var err error
	c.Conn, err = net.Dial("tcp",
		fmt.Sprintf("%s:%d", c.Service.Host, c.Service.Port))
	defer c.Conn.Close()

	if err != nil {
		fmt.Printf("conn server failed, err： %v\n", err)
		return
	}

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = c.Conn.Read(message.Buffer)
		if err != nil {
			panic(err)
		}
		message.Name = "in-net"
		c.Messages <- message
	}
}

func (c *Client) Run() {
	go c.TcpClient()

	go c.TunServer()

	for {
		select {
		case message := <-c.Messages:
			switch message.Name {
			case "in-tun":
				c.Conn.Write(message.Buffer[:message.Len])
			case "in-net":
				c.Tun.Write(message.Buffer[:message.Len])
			}
		}
	}
}
