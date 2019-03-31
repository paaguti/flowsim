package tcp

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"errors"
	"net"
	// "io"
	// "os"

	"crypto/rand"
	"time"

	common "github.com/paaguti/flowsim/common"
)

func matcher(cmd string) (string, string, string, error) {
	expr := regexp.MustCompile(`GET (\d+)/(\d+) (\d+)`)
	parsed := expr.FindStringSubmatch(cmd)
	if len(parsed) == 4 {
		return parsed[1], parsed[2], parsed[3], nil
	}
	return "", "", "", errors.New(fmt.Sprintf("Unexpected request %s", cmd))
}

func handleConn(conn *net.TCPConn) {
	var run, total, bunch string

	defer conn.Close()

	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if common.FatalError(err) != nil {
			return
		}
		// output message received
		fmt.Printf("Message received at %s: %s", time.Now().Format(time.RFC3339), string(message))

		// Checked in the client
		run, total, bunch, err = matcher(strings.ToUpper(string(message)))
		if common.FatalError(err) != nil {
			continue
		}

		run_iter, _ := strconv.Atoi(run)
		total_iter, _ := strconv.Atoi(total)
		bunch_len, _ := strconv.Atoi(bunch)

		// conn.Write([]byte(fmt.Sprintf("run %d of %d... should send %d bytes\n",run_iter, total_iter, bunch_len)))

		testBunch := make([]byte, bunch_len)
		// numRead, err := io.ReadFull(zero, testBunch)
		numRead, err := rand.Read(testBunch)
		// fmt.Printf("Read %d bytes from /dev/zero\n",len(testBunch))
		if common.FatalError(err) != nil {
			continue
		}
		fmt.Printf("Sending %d bytes...\n", numRead)
		conn.Write(testBunch)
		if run_iter == total_iter {
			// fmt.Println("This should kill this TCP server thread")
			break
		}
	}
	fmt.Println("Connection closed...")
}

func Server(ip string, port int, single bool, tos int) {

	listenAddrStr := net.JoinHostPort(ip, strconv.Itoa(port))

	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddrStr)
	if common.FatalErrorf(err, "Error resolving %s:%d (%v)\n", ip, port, err) != nil {
		return
	}

	ln, err := net.ListenTCP("tcp", listenAddr)
	if common.FatalErrorf(err, "Error binding server to %s\n", listenAddr) != nil {
		return
	}

	fmt.Printf("Listening at %s\n", listenAddr)
	for {
		// accept connection on port
		conn, err := ln.AcceptTCP()
		if common.FatalErrorln(err, "Error accepting connection") != nil {
			continue
		}

		err = common.SetTcpTos(conn, tos)
		if common.FatalErrorln(err, "Error setting TOS") != nil {
			continue
		}

		if single {
			handleConn(conn)
			break
		} else {
			go handleConn(conn)
		}
	}
}
