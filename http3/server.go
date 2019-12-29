package http3

import (
	"context"
	// "crypto/tls"
	// "crypto/x509"
	"fmt"
	quic "github.com/lucas-clemente/quic-go"
	http3 "github.com/lucas-clemente/quic-go/http3"
	common "github.com/paaguti/flowsim/common"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	// "time"
)

func getInt(query url.Values, field string) (int, bool) {
	params, ok := query[field]
	if !ok {
		log.Println("Url Param '" + field + "' is missing")
		return -1, false
	}
	value, err := strconv.Atoi(params[0])
	return value, err == nil
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func Server(ip string, port int, single bool, tos int, certs string) {

	srvClosed := make(chan int)
	mux := http.NewServeMux()
	mux.HandleFunc("/flowsim/request", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		log.Printf("Received %s%v", r.URL.Path, query)

		requested, ok := getInt(query, "bytes")
		if !ok {
			log.Println("Couldn't decode " + r.URL.Path)
			return
		}
		pass, ok := getInt(query, "pass")
		if !ok {
			log.Println("Couldn't decode " + r.URL.Path)
			return
		}
		total, ok := getInt(query, "of")
		if !ok {
			log.Println("Couldn't decode " + r.URL.Path)
			return
		}

		// fmt.Fprintf(w, "So you are requesting "+bytes+" bytes in pass "+pass+" of "+of+" from me...")
		fmt.Fprintln(w, common.RandStringBytes(requested))
		if pass == total {
		}
		log.Printf("Served %d bytes in pass %d of %d", requested, pass, total)
	})

	mux.HandleFunc("/flowsim/close", func(w http.ResponseWriter, req *http.Request) {
		// if req.URL.Path != "/" {
		// 	http.NotFound(w, req)
		// 	log.Fatal("Can't serve " + req.URL.Path)
		// } else {
		fmt.Fprintf(w, "Bye!\n")
		// }
		if single {
			log.Printf("And here we should stop")
			srvClosed <- 1
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			log.Fatal("Can't serve " + req.URL.Path)
		} else {
			fmt.Fprintf(w, "Welcome to the flowsim QUIC server!")
		}
	})

	for {
		var server http3.Server
		// go func() {
		//   log.Println(http.ListenAndServe("localhost:6060",nil))
		// }()

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
