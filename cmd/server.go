package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"net"
)

type Server struct {
	DeviceName string
	Listen string
	Port int
}

func (s *Server) Init() {
	s.Listen = "0.0.0.0"
	s.Port = 5050
	s.DeviceName = "tun0"
}

func (s *Server) Process(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		var buf [2048]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed : %v\n", err)
			break
		}


	}
}

func(s *Server) Run() {
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
		go s.Process(conn)
	}

	go Mktun()
}

func ServerRun(cmd *cobra.Command, args []string) {
	server := &Server{}

	server.Init()

	server.Run()
}


// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `server`,
	Run: ServerRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
