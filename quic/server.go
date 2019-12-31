package quic

import (
	"bufio"
	"bytes"
	"context"
	// "crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	common "github.com/paaguti/flowsim/common"
)

// Start a server that echos all data on the first stream opened by the client
func Server(ip string, port int, single bool, dscp int) error {

	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if common.FatalError(err) != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if common.FatalError(err) != nil {
		return err
	}

	err = common.SetUdpTos(conn, dscp)
	if common.FatalError(err) != nil {
		return err
	}

	tlsConfig, err := common.ServerTLSConfig("/etc", "flowsim-quic")
	if common.FatalError(err) != nil {
		return err
	}

	listener, err := quic.Listen(conn, tlsConfig, nil)
	defer listener.Close()
	if common.FatalError(err) != nil {
		return err
	}

	for {

		sess, err := listener.Accept(context.Background())

		if common.FatalError(err) != nil {
			return err
		}
		if single {
			err = quicHandler(sess)
			time.Sleep(500 * time.Millisecond)
			return err
		}
		go quicHandler(sess)
	}
}

func quicHandler(sess quic.Session) error {

	log.Println("Entering quicHandler")

	stream, err := sess.AcceptStream(context.Background())

	if common.FatalError(err) != nil {
		return err
	}
	// defer stream.Close()
	log.Println("Got a stream")
	msgbuf := make([]byte, 128)
	reader := bufio.NewReaderSize(stream, 128)
	for end := false; !end; {
		log.Println("In server loop")
		n, err := reader.Read(msgbuf)
		common.FatalError(err)
		if err != nil {
			if end == true {
				log.Println("Bye!")
				return nil
			}
			return err
		}

		log.Printf("In server loop: got %d bytes: %s", n, msgbuf)
		wbuf, _end, err := parseCmd(string(msgbuf))
		if common.FatalError(err) != nil {
			return err
		}
		end = _end
		_, err = io.Copy(stream, bytes.NewBuffer(wbuf))
		if common.FatalError(err) == nil {
			log.Println("Sent bytes")
		}
	}
	time.Sleep(1 * time.Second)
	return nil
}

// From flowsim TCP
func matcher(cmd string) (string, string, string, error) {
	expr := regexp.MustCompile(`GET (\d+)/(\d+) (\d+)`)
	parsed := expr.FindStringSubmatch(cmd)
	if len(parsed) == 4 {
		return parsed[1], parsed[2], parsed[3], nil
	}
	return "", "", "", errors.New(fmt.Sprintf("Unexpected request %s", cmd))
}

/*
* Purpuse: parse get Command from client
*         and generate a buffer with random bytes
* Return: byte buffer to send or nil on error
*         boolean: true id last bunch
*         error or nil if all went well
*
* Uses crypto/rand, which is already imported for key handling
 */
func parseCmd(strb string) ([]byte, bool, error) {
	// log.Printf("Server: Got %s", strb)
	iter, total, bunchStr, err := matcher(strb)
	if err == nil {
		bunch, _ := strconv.Atoi(bunchStr) // ignore error, wouldn't have parsed the command
		nb := common.RandBytes(bunch)
		if err != nil {
			log.Printf("ERROR while filling random buffer: %v\n", err)
			return nil, iter == total, err
		}
		log.Printf("Sending %d bytes\n", len(nb))
		return nb, iter == total, err
	}
	return nil, false, err
}
