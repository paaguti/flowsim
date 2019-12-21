package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	common "github.com/paaguti/flowsim/common"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"
)

type Transfer struct {
	XferStart string
	XferTime  string
	XferBytes int
	XferIter  int
}

type Result struct {
	Protocol string
	Server   string
	Burst    int
	Start    string
	Times    []Transfer
}

func mkTransfer(serverAddr string, iter int, total int, tsize int, dscp int, t time.Time, tlsConfig *tls.Config) *Transfer {
	var tr *http.Transport
	var proto string

	tr = &http.Transport{
		MaxIdleConns:       1,
		IdleConnTimeout:    1 * time.Second,
		DisableCompression: true,
		DisableKeepAlives:  true,
		TLSClientConfig:    tlsConfig,
	}

	// log.Printf("Using TLS config: %v", tlsConfig.RootCAs)

	if tlsConfig.RootCAs == nil {
		proto = "http"
	} else {
		proto = "https"
	}
	log.Printf("Starting an %s client", proto)
	//
	// Until you get the DSCP right, just make it part of the request
	//
	server_url := fmt.Sprintf("%s://%s/flowsim/request?bytes=%d&pass=%d&of=%d&dscp=%d", proto, serverAddr, tsize, iter, total, dscp)
	//
	//
	//
	req, err := http.NewRequest("GET", server_url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Transport: tr}
	// defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	fmt.Printf("Got %d bytes back\n", len(body))
	// if len(body) != bunch {
	// 	log.Fatal(string(body))
	// }

	return &Transfer{
		XferStart: t.Format(time.RFC3339),
		XferTime:  time.Since(t).String(),
		XferBytes: len(body),
		XferIter:  iter,
	}
}

func Client(ip string, port int, iter int, interval int, bunch int, dscp int, certs string) error {

	var resultProto string
	var tlsConfig *tls.Config

	if certs != "" {
		log.Printf("Trying to read certificates from %s", certs)

		// path.Join(path.Dir(filename), "data.csv")
		caCert, err := ioutil.ReadFile(path.Join(certs, "server.crt"))
		if err != nil {
			return common.FatalError(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		cert, err := tls.LoadX509KeyPair(path.Join(certs, "client.crt"), path.Join(certs, "client.key"))
		if err != nil {
			return common.FatalError(err)
		}

		tlsConfig = &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		}
		resultProto = "HTTPS"
	} else {
		tlsConfig = &tls.Config{}
		resultProto = "HTTP"
	}
	serverAddrStr := net.JoinHostPort(ip, strconv.Itoa(port))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := Result{
		Protocol: resultProto,
		Server:   serverAddrStr,
		Burst:    bunch,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	result.Times[0] = *mkTransfer(serverAddrStr, 1, iter, bunch, dscp, time.Now(), tlsConfig)
	currIter := 1

	if iter > 1 {
		done := make(chan bool, 1)
		for {
			select {
			case t := <-ticker.C:
				currIter++
				if currIter >= iter {
					close(done)
				}
				result.Times[currIter-1] = *mkTransfer(serverAddrStr, currIter, iter, bunch, dscp, t, tlsConfig)
			case <-done:
				// fmt.Fprintf(os.Stderr, "Finished...\n\n")
				common.PrintJSon(result)
				return nil
			}
		}
	}
	return nil
}
