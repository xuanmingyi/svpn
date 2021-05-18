package cmd

import (
	"github.com/spf13/cobra"
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

func(s *Server) Run() {
	
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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
