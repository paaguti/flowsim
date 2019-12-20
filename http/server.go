package http

import (
	"context"
	"fmt"
	common "github.com/paaguti/flowsim/common"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

func Server(ip string, port int, single bool, tos int) {

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
		fmt.Fprintf(w, common.RandStringBytes(requested))
		if pass == total {
			log.Printf("And here we should stop")
			srvClosed <- 1
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			log.Fatal("Can't serve " + req.URL.Path)
		} else {
			fmt.Fprintf(w, "Welcome to the flowsim HTTP server!")
		}
	})
	var srv *http.Server

	srv = &http.Server{
		Addr:           net.JoinHostPort(ip, strconv.Itoa(port)),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	for {
		go func() {
			srv.ListenAndServe()
		}()
		<-srvClosed
		if single {
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
			log.Printf("HTTP server shutdown")
			break
		}
	}

}
