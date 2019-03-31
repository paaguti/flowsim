package quic

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"crypto/tls"
	"math/rand"

	quic "github.com/lucas-clemente/quic-go"
	common "github.com/paaguti/flowsim/common"
)

type Transfer struct {
	XferTime  time.Duration
	XferBytes int
	XferIter  int
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

	session, err := quic.Dial(udpConn, updAddr, addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if common.FatalError(err) != nil {
		return err
	}
	defer session.Close()

	// fmt.Printf("Opened session for %s\n", addr)
	buf := make([]byte, bunch)
	stream, err := session.OpenStreamSync()
	if common.FatalError(err) != nil {
		return err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	fmt.Printf("{\n  \"burst\" : \"%d\",\n  \"server\" : \"%s\",\n", bunch, addr)
	fmt.Printf("  \"start\" : \"%s\",\n", time.Now().Format(time.RFC3339))
	fmt.Printf("  \"times\": [\n")

	times := make([]Transfer, iter)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	times[0] = *mkTransfer(stream, buf, 1, iter, time.Now())
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
				times[currIter-1] = *mkTransfer(stream, buf, currIter, iter, t)
			case <-done:
				for i := range times {
					sep := ','
					if i == len(times)-1 {
						sep = ' '
					}
					fmt.Printf("    {\"time\" : \"%v\", \"xferd\" : \"%d\", \"n\" : \"%d\"",
						times[i].XferTime, times[i].XferBytes, times[i].XferIter)
					fmt.Printf("}%c \n", sep)
				}
				fmt.Printf("  ]\n}\n")
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
		XferTime:  time.Since(t),
		XferBytes: n,
		XferIter:  current,
	}
}
