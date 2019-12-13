package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	//	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flowsim",
	Short: "A TCP/QUIC server/client to simulate ABR traffic",
	Long: `A TCP/QUIC server/client to simulate ABR traffic in one app.
Follows the iperf3 way of life integrating both server and client

This version is based on an older fork of the lucas-clemente/quic-go code
that supports the spin bit`,
}

func Execute() {
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
