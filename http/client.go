package http

import (
	"fmt"
	common "github.com/paaguti/flowsim/common"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
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

func mkTransfer(serverAddr string, iter int, total int, tsize int, dscp int, t time.Time) *Transfer {
	//
	// Until you get the DSCP right, just make it part of the request
	//
	server_url := fmt.Sprintf("http://%s/flowsim/request?bytes=%d&pass=%d&of=%d&dscp=%d", serverAddr, tsize, iter, total, dscp)

	tr := &http.Transport{
		MaxIdleConns:       1,
		IdleConnTimeout:    1 * time.Second,
		DisableCompression: true,
		DisableKeepAlives:  true,
	}

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

func Client(ip string, port int, iter int, interval int, bunch int, dscp int) error {

	serverAddrStr := net.JoinHostPort(ip, strconv.Itoa(port))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := Result{
		Protocol: "HTTP",
		Server:   serverAddrStr,
		Burst:    bunch,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	result.Times[0] = *mkTransfer(serverAddrStr, 1, iter, bunch, dscp, time.Now())
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
				result.Times[currIter-1] = *mkTransfer(serverAddrStr, currIter, iter, bunch, dscp, t)
			case <-done:
				// fmt.Fprintf(os.Stderr, "Finished...\n\n")
				common.PrintJSon(result)
				return nil
			}
		}
	}
	return nil
}
