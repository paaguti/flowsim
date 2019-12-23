package cmd

import (
	"fmt"
	common "github.com/paaguti/flowsim/common"
	"github.com/paaguti/flowsim/http"
	"github.com/paaguti/flowsim/quic"
	"github.com/paaguti/flowsim/tcp"
	"github.com/spf13/cobra"
	"path"
)

var clientIp string
var clientPort int
var clientInterval int
var clientIter int
var clientTos string
var clientBurst string
var clientQuic bool
var clientHttp bool
var clientHttps bool
var clientIpv6 bool

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a flowsim TCP, HTTP or QUIC client",
	Long: `Will run flowsim in client mode
and try to talk to a flowsim server.
CAVEAT: flowsim server and client mode have be in the same mode.`,

	Run: func(cmd *cobra.Command, args []string) {

		tos, err := Dscp(clientTos)
		if err != nil {
			fmt.Printf("Warning: %v, TOS will be %d instead of %s \n", err, tos, clientTos)
		}
		burstSize, err := utoi(clientBurst)
		if err != nil {
			fmt.Printf("Warning: %v, generating %d byte bursts\n", err, burstSize)
		}
		useIp, err := common.FirstIP(clientIp, clientIpv6)
		common.FatalError(err)
		if clientQuic {
			quic.Client(useIp, clientPort, clientIter, clientInterval, burstSize, tos*4)
		} else if clientHttp || clientHttps {
			certDir := ""
			if clientHttps {
				if exePath, err := common.ExePath(); err == nil {
					certDir = path.Dir(exePath)
				}
			}
			http.Client(useIp, clientPort, clientIter, clientInterval, burstSize, tos*4, certDir)
		} else {
			tcp.Client(useIp, clientPort, clientIter, clientInterval, burstSize, tos*4)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringVarP(&clientIp, "ip", "I", "localhost", "IP address or host name of the flowsim server to talk to")
	clientCmd.PersistentFlags().IntVarP(&clientPort, "port", "p", 8081, "TCP port of the flowsim server")
	clientCmd.PersistentFlags().IntVarP(&clientIter, "iter", "n", 6, "Number of bursts")
	clientCmd.PersistentFlags().IntVarP(&clientInterval, "interval", "t", 10, "Interval in secs between bursts")
	clientCmd.PersistentFlags().StringVarP(&clientBurst, "burst", "N", "1M", "Size of each burst (as x(.xxx)?[kmgtKMGT]?)")
	clientCmd.PersistentFlags().StringVarP(&clientTos, "TOS", "T", "CS0", "Value of the DSCP field in the IP packets (valid int or DSCP-Id)")
	clientCmd.PersistentFlags().BoolVarP(&clientQuic, "quic", "Q", false, "Use QUIC (default is TCP)")
	clientCmd.PersistentFlags().BoolVarP(&clientHttp, "http", "H", false, "Use HTTP (default is TCP)")
	clientCmd.PersistentFlags().BoolVarP(&clientHttps, "https", "S", false, "Use HTTPS (default is TCP)")
	clientCmd.PersistentFlags().BoolVarP(&clientIpv6, "ipv6", "6", false, "Use IPv6 (default is IPv4)")
}
