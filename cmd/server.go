package cmd

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/spf13/cobra"
)

type Server struct {
	DeviceName string
	Listen     string
	Port       int

	File     *os.File
	Messages chan Message
	Connects []net.Conn
}

func (s *Server) Init() {
	s.Connects = make([]net.Conn, 0)
	s.Messages = make(chan Message)
}

func (s *Server) Process(conn net.Conn) {
	var err error
	defer conn.Close()

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = conn.Read(message.Buffer)
		if err != nil {
			fmt.Printf("read from conn failed : %v\n", err)
			break
		}
		message.Name = "in-net"
		s.Messages <- message

	}
}

func (s *Server) Mktun() {
	var err error

	s.File, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	ifr := make([]byte, 18)
	copy(ifr, []byte(s.DeviceName))
	ifr[16] = 0x01
	ifr[17] = 0x10

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.File.Fd()), uintptr(0x400454ca),
		uintptr(unsafe.Pointer(&ifr[0])))

	if errno != 0 {
		panic(errno)
	}

	for {
		var message Message
		message.Buffer = make([]byte, 2048)
		message.Len, err = s.File.Read(message.Buffer)
		if err != nil {
			panic(err)
		}
		message.Name = "in-tun"
		s.Messages <- message
	}

}

func (s *Server) TcpServer() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Listen, s.Port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed: err： %v\n", err)
			continue
		}
		s.Connects = append(s.Connects, conn)
		go s.Process(conn)
	}
}

func (s *Server) Run() {
	// 监听网卡
	go s.Mktun()

	// 监听网络
	go s.TcpServer()

	for {
		select {
		case message := <-s.Messages:
			if message.Name == "in-tun" {
				fmt.Println("read from tun")
				for _, client := range s.Connects {
					client.Write(message.Buffer[:message.Len])
				}
			} else if message.Name == "in-net" {
				s.File.Write(message.Buffer[:message.Len])
				fmt.Println("read from net")
			}
		}
	}

}

func ServerRun(cmd *cobra.Command, args []string) {
	server := &Server{
		Listen:     "0.0.0.0",
		Port:       5050,
		DeviceName: "tun0",
	}

	server.Init()

	server.Run()
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  `server`,
	Run:   ServerRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
