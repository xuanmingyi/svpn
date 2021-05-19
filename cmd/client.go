package cmd

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/spf13/cobra"
)

type Client struct {
	DeviceName string
	Host       string
	Port       int

	File     *os.File
	Messages chan Message
	Conn     net.Conn
}

func (c *Client) Init() {
	c.Messages = make(chan Message)
}

func (c *Client) Mktun() {
	var err error

	c.File, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	ifr := make([]byte, 18)
	copy(ifr, []byte(c.DeviceName))
	ifr[16] = 0x01
	ifr[17] = 0x10

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.File.Fd()), uintptr(0x400454ca),
		uintptr(unsafe.Pointer(&ifr[0])))

	if errno != 0 {
		panic(errno)
	}

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = c.File.Read(message.Buffer)
		if err != nil {
			panic(err)
		}
		message.Name = "in-tun"
		c.Messages <- message
	}

}

func (c *Client) TcpClient() {
	var err error
	c.Conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	defer c.Conn.Close()

	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
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
	// 监听网卡
	go c.Mktun()

	// 监听网络
	go c.TcpClient()

	for {
		select {
		case message := <-c.Messages:
			if message.Name == "in-tun" {
				c.Conn.Write(message.Buffer[:message.Len])
				fmt.Println("read from tun sent over netowrk")

			} else if message.Name == "in-net" {
				c.File.Write(message.Buffer[:message.Len])
				fmt.Println("read from net")
			}
		}
	}

}

func ClientRun(cmd *cobra.Command, args []string) {
	client := &Client{
		Host:       "18.1.1.101",
		Port:       5050,
		DeviceName: "tun0",
	}

	client.Init()

	client.Run()
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client",
	Long:  `client`,
	Run:   ClientRun,
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
