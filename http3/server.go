package http3

import (
	"context"
	// "crypto/tls"
	// "crypto/x509"
	// "fmt"
	quic "github.com/lucas-clemente/quic-go"
	http3 "github.com/lucas-clemente/quic-go/http3"
	common "github.com/paaguti/flowsim/common"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	// "net/url"
	"path"
	"strconv"
	// "time"
)

// func getInt(query url.Values, field string) (int, bool) {
// 	params, ok := query[field]
// 	if !ok {
// 		log.Println("Url Param '" + field + "' is missing")
// 		return -1, false
// 	}
// 	value, err := strconv.Atoi(params[0])
// 	return value, err == nil
// }

// type apiHandler struct{}

// func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func Server(ip string, port int, single bool, tos int, certs string) {

	srvClosed := make(chan int)
	mux := common.MakeHTTPMux(single, srvClosed)
	for {
		var server http3.Server

		go func() {
			bCap := net.JoinHostPort(ip, strconv.Itoa(port))
			log.Printf("Starting HTTP3 server at %s", bCap)
			server = http3.Server{
				Server: &http.Server{
					Handler: mux,
					Addr:    bCap,
				},
				QuicConfig: &quic.Config{},
			}
			certFile := path.Join(certs, "flowsim-server.crt")
			keyFile := path.Join(certs, "flowsim-server.key")
			err := server.ListenAndServeTLS(certFile, keyFile)
			if err != nil {
				log.Printf(" From http3 server: %v", err)
			}
		}()
		<-srvClosed
		if single {
			server.Shutdown(context.Background())
			log.Printf("http3 server shutdown")
			break
		}
	}
}
