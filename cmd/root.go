package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flowsim",
	Short: "A UDP/TCP/HTTP/QUIC server/client to simulate ABR traffic",
	Long: `A UDP/TCP/HTTP/QUIC server/client to simulate ABR traffic in one app.
Follows the iperf3 way of life integrating both server and client

This version is based on an older fork of the lucas-clemente/quic-go code
that supports the spin bit`,
}

func Execute() {
	// if exepath, err := ExePath(); err != nil {
	// 	fmt.Println("Couldn't get path to executable!")
	// } else {
	// 	fmt.Println("Executing ", exepath)
	// }
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
