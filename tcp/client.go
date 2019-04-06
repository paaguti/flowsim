package tcp

import (
	"fmt"
	"net"
	// "log"
	common "github.com/paaguti/flowsim/common"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type Transfer struct {
	XferStart string
	XferTime  string
	XferBytes int
	XferIter  int
}

func mkTransfer(conn *net.TCPConn, iter int, total int, tsize int, t time.Time) *Transfer {
	// fmt.Fprintf(os.Stderr, "Launching at %v\n", t)
	// send to socket
	fmt.Fprintf(conn, fmt.Sprintf("GET %d/%d %d\n", iter, total, tsize))
	// listen for reply
	readBuffer := make([]byte, tsize)
	// fmt.Fprintf(os.Stderr, "Trying to read %d bytes back...", len(readBuffer))
	readBytes, err := io.ReadFull(conn, readBuffer)
	common.FatalError(err)

	// fmt.Fprintf(os.Stderr, "Effectively read %d bytes\n", readBytes)
	return &Transfer{
		XferStart: t.Format(time.RFC3339),
		XferTime:  time.Since(t).String(),
		XferBytes: readBytes,
		XferIter:  iter,
	}
}

type Result struct {
	Protocol string
	Server   string
	Burst    int
	Start    string
	Times    []Transfer
}

func Client(host string, port int, iter int, interval int, burst int, tos int) {

	serverAddrStr := net.JoinHostPort(host, strconv.Itoa(port))

	server, err := net.ResolveTCPAddr("tcp", serverAddrStr)
	if common.FatalErrorf(err, "Error resolving %s\n", serverAddrStr) != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, server)
	if common.FatalErrorf(err, "Error connecting to %s: %v\n", serverAddrStr) != nil {
		return
	}
	defer conn.Close()
	// fmt.Fprintf(os.Stderr, "Talking to %s\n", serverAddrStr)

	err = common.SetTcpTos(conn, tos)
	common.FatalError(err)

	// fmt.Fprintf(os.Stderr,("Starting at  %v\n", time.Now())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := Result{
		Protocol: "TCP",
		Server:   serverAddrStr,
		Burst:    burst,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	result.Times[0] = *mkTransfer(conn, 1, iter, burst, time.Now())
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
				result.Times[currIter-1] = *mkTransfer(conn, currIter, iter, burst, t)
			case <-done:
				// fmt.Fprintf(os.Stderr, "Finished...\n\n")
				common.PrintJSon(result)
			}
		}
	}
}
