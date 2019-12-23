package http

import (
	"crypto/tls"
	// "crypto/x509"
	"fmt"
	common "github.com/paaguti/flowsim/common"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	// "path"
	"strconv"
	"strings"
	"time"
)

func mkTransfer(serverAddr string, iter int, total int, tsize int, dscp int, t time.Time, tlsConfig *tls.Config) (*common.Transfer, string) {
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

	if common.IsSecureConfig(tlsConfig) {
		proto = "https"
	} else {
		proto = "http"
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

	return &common.Transfer{
		XferStart: t.Format(time.RFC3339),
		XferTime:  time.Since(t).String(),
		XferBytes: len(body),
		XferIter:  iter,
	}, proto
}

func Client(ip string, port int, iter int, interval int, bunch int, dscp int, certs string) error {

	var resultProto string
	var tlsConfig *tls.Config

	if certs == "" {
		tlsConfig = &tls.Config{}
	} else {
		log.Printf("Trying to read certificates from %s", certs)
		var err error
		// tlsConfig, err = common.HttpsClientTLSConfig(certs)
		tlsConfig, err = common.ClientTLSConfig(certs)
		if err != nil {
			return err
		}
	}

	serverAddrStr := net.JoinHostPort(ip, strconv.Itoa(port))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := common.Result{
		Protocol: "",
		Server:   serverAddrStr,
		Burst:    bunch,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]common.Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	measure, resultProto := mkTransfer(serverAddrStr, 1, iter, bunch, dscp, time.Now(), tlsConfig)
	result.Times[0] = *measure
	result.Protocol = strings.ToUpper(resultProto)
	currIter := 1

	if iter > 1 {
		done := make(chan bool, 1)
		defer close(done)
		for {
			select {
			case t := <-ticker.C:
				currIter++
				measure, _ = mkTransfer(serverAddrStr, currIter, iter, bunch, dscp, t, tlsConfig)
				result.Times[currIter-1] = *measure
				if currIter >= iter {
					done <- true
				}
			case <-done:
				// fmt.Fprintf(os.Stderr, "Finished...\n\n")
				common.PrintJSon(result)
				return nil
			}
		}
	}
	return nil
}
