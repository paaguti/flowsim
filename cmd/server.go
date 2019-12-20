package cmd

import (
	"fmt"
	common "github.com/paaguti/flowsim/common"
	"github.com/paaguti/flowsim/http"
	"github.com/paaguti/flowsim/quic"
	"github.com/paaguti/flowsim/tcp"
	"github.com/spf13/cobra"
)

var serverIp string
var serverPort int
var serverSingle bool
var serverTos string
var serverQuic bool
var serverHttp bool
var serverIpv6 bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a flowsim TCP or QUIC server",
	Long: `Start an TCP or QUIC ABR server.
It will basically sit there and wait for the client to request bunches of data
over a TCP, HTTP or QUIC connection
Payload is filled with random bytes`,
	Run: func(cmd *cobra.Command, args []string) {
		tos, err := Dscp(serverTos)
		if err != nil {
			fmt.Printf("Warning: %v, TOS will be %d instead of %s \n", err, tos, serverTos)
		}
		useIp, err := common.FirstIP(serverIp, serverIpv6)
		common.FatalError(err)

		if serverQuic {
			// fmt.Println("Warning: QUIC doesn't support setting DSCP yet!")
			quic.Server(useIp, serverPort, serverSingle, tos*4)
		} else if serverHttp {
			http.Server(useIp, serverPort, serverSingle, tos*4)
		} else {
			tcp.Server(useIp, serverPort, serverSingle, tos*4)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&serverIp, "ip", "I", "localhost", "IP address or host name bound to the flowsim server")
	serverCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 8081, "TCP port bound to the flowsim server")
	serverCmd.PersistentFlags().BoolVarP(&serverSingle, "one-off", "1", false, "Just accept one connection and quit (default is run until killed)")
	serverCmd.PersistentFlags().StringVarP(&serverTos, "TOS", "T", "CS0", "Value of the DSCP field in the IP layer (number or DSCP id)")
	serverCmd.PersistentFlags().BoolVarP(&serverQuic, "quic", "Q", false, "Use QUIC (default is TCP)")
	serverCmd.PersistentFlags().BoolVarP(&serverHttp, "http", "H", false, "Use HTTP (default is TCP)")
	serverCmd.PersistentFlags().BoolVarP(&serverIpv6, "ipv6", "6", false, "Use IPv6 (default is IPv4)")
}
