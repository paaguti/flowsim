package http3

import (
	"fmt"
	common "github.com/paaguti/flowsim/common"
	quic "github.com/paaguti/quic-go"
	h2quic "github.com/paaguti/quic-go/h2quic"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func mkTransfer(serverAddr string, iter int, total int, tsize int, dscp int, t time.Time, roundTripper *h2quic.RoundTripper) (*common.Transfer, string) {

	//
	// Always use https:// in the URL
	// Until you get the DSCP right, just make it part of the request
	//
	server_url := fmt.Sprintf("https://%s/flowsim/request?bytes=%d&pass=%d&of=%d&dscp=%d", serverAddr, tsize, iter, total, dscp)
	//
	//
	//
	log.Printf("Starting an h2quic client to\n%s", server_url)

	req, err := http.NewRequest("GET", server_url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{
		Transport: roundTripper,
	}

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
	}, ""
}

func Client(ip string, port int, iter int, interval int, bunch int, dscp int, certs string) error {

	// log.Println("H2QUIC client test...")
	versions := quic.SupportedVersions
	versions = append([]quic.VersionNumber{quic.VersionTLS}, versions...)

	roundTripper := &h2quic.RoundTripper{
		QuicConfig: &quic.Config{Versions: versions},
	}
	defer roundTripper.Close()

	serverAddrStr := net.JoinHostPort(ip, strconv.Itoa(port))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := common.Result{
		Protocol: "H2QUIC",
		Server:   serverAddrStr,
		Burst:    bunch,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]common.Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	measure, _ := mkTransfer(serverAddrStr, 1, iter, bunch, dscp, time.Now(), roundTripper)
	result.Times[0] = *measure

	currIter := 1
	if iter > 1 {
		done := make(chan bool, 1)
		defer close(done)
		for {
			select {
			case t := <-ticker.C:
				currIter++
				measure, _ = mkTransfer(serverAddrStr, currIter, iter, bunch, dscp, t, roundTripper)
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
	common.PrintJSon(result)
	return nil
}
