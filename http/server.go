package http

import (
	"context"
	"crypto/tls"
	// "crypto/x509"
	// "fmt"
	common "github.com/paaguti/flowsim/common"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	// "net/url"
	"path"
	"strconv"
	"time"
)

func Server(ip string, port int, single bool, tos int, certs string) {

	srvClosed := make(chan int)
	mux := common.MakeHTTPMux(single, srvClosed)

	var srv *http.Server

	if certs == "" {
		srv = &http.Server{
			Addr:           net.JoinHostPort(ip, strconv.Itoa(port)),
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      &tls.Config{},
		}
	} else {
		tlsConfig, err := common.ServerTLSConfig(certs)
		// tlsConfig, err := common.HttpsServerTLSConfig(certs)
		if err != nil {
			return
		}
		srv = &http.Server{
			Addr:           net.JoinHostPort(ip, strconv.Itoa(port)),
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      tlsConfig,
		}
	}

	for {
		go func() {
			if certs == "" {
				log.Println("Starting HTTP server")
				srv.ListenAndServe()
			} else {
				log.Println("Starting HTTPS server")
				srv.ListenAndServeTLS(path.Join(certs, "flowsim-server.crt"), path.Join(certs, "flowsim-server.key"))
				// srv.ListenAndServeTLS("", "")
			}
		}()
		<-srvClosed
		if single {
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
			if certs == "" {
				log.Println("HTTP server shutdown")
			} else {
				log.Println("HTTPS server shutdown")
			}
			break
		}
	}
}
