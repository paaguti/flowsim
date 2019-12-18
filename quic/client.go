package quic

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	// "context"
	"crypto/tls"
	"math/rand"

	common "github.com/paaguti/flowsim/common"
	//
	// use the fork with the Spinbit and VEC implementation
	// I have forked ferrieux/quic-go to keep a stable version
	//
	// quic "github.com/ferrieux/quic-go"
	quic "github.com/paaguti/quic-go"
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

func Client(ip string, port int, iter int, interval int, bunch int, dscp int) error {

	udpFamily, err := common.UdpFamily(ip)
	if common.FatalError(err) != nil {
		return err
	}

	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	updAddr, err := net.ResolveUDPAddr(udpFamily, addr)
	if common.FatalError(err) != nil {
		return err
	}

	udpConn, err := net.ListenUDP(udpFamily, &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if common.FatalError(err) != nil {
		return err
	}

	err = common.SetUdpTos(udpConn, dscp)
	if common.FatalError(err) != nil {
		return err
	}

	var config *quic.Config

	config = &quic.Config{Versions: []quic.VersionNumber{quic.VersionGQUIC39}}

	// config := quic.PopulateClientConfig(nil, false)

	session, err := quic.Dial(udpConn, updAddr, addr, &tls.Config{InsecureSkipVerify: true},
		config)
	if common.FatalError(err) != nil {
		return err
	}
	defer session.Close(err)

	// fmt.Printf("Opened session for %s\n", addr)
	buf := make([]byte, bunch)
	// This is for the latest version of quic-go
	// stream, err := session.OpenStreamSync(context.Background())
	//
	// revert to get the spin bit running
	//
	stream, err := session.OpenStreamSync()
	if common.FatalError(err) != nil {
		return err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := Result{
		Protocol: "QUIC",
		Server:   addr,
		Burst:    bunch,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	result.Times[0] = *mkTransfer(stream, buf, 1, iter, time.Now())
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
				result.Times[currIter-1] = *mkTransfer(stream, buf, currIter, iter, t)
			case <-done:
				common.PrintJSon(result)
				// fmt.Printf("\nFinished...\n\n")
				return nil
			}
		}
	}
	return nil
}

func mkTransfer(stream quic.Stream, buf []byte, current int, iter int, t time.Time) *Transfer {

	message := fmt.Sprintf("GET %d/%d %d\n", current, iter, len(buf))
	// fmt.Printf("Client: (%v) Sending > %s", t, message)

	_, err := stream.Write([]byte(message))
	if common.FatalError(err) != nil {
		return nil
	}

	n, err := io.ReadFull(stream, buf)
	common.FatalError(err)
	// fmt.Printf("Client: Got %d bytes back\n", n)
	return &Transfer{
		XferStart: t.Format(time.RFC3339),
		XferTime:  time.Since(t).String(),
		XferBytes: n,
		XferIter:  current,
	}
}
