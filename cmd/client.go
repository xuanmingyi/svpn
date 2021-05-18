package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
)

func ClientRun(cmd *cobra.Command, args []string) {
	conn, err := net.Dial("tcp", "127.0.0.1:5050")
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return
	}

	conn.Write([]byte("sss"))
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `client`,
	Run: ClientRun,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
