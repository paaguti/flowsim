package quic

import (
	"context"
	"fmt"
	quic "github.com/lucas-clemente/quic-go"
	common "github.com/paaguti/flowsim/common"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

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

	// config = &quic.Config{Versions: []quic.VersionNumber{quic.VersionGQUIC39}}
	config = &quic.Config{}

	// TODO: include certificate configuration for a better TLS verification

	tlsConfig, err := common.ClientTLSConfig("")
	tlsConfig.NextProtos = []string{"flowsim-quic"}

	if common.FatalError(err) != nil {
		return err
	}

	session, err := quic.Dial(udpConn, updAddr, addr, tlsConfig, config)
	if common.FatalError(err) != nil {
		return err
	}
	defer session.Close()

	// fmt.Printf("Opened session for %s\n", addr)
	buf := make([]byte, bunch)

	stream, err := session.OpenStreamSync(context.Background())

	if common.FatalError(err) != nil {
		return err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initWait := r.Intn(interval*50) / 100.0
	time.Sleep(time.Duration(initWait) * time.Second)

	result := common.Result{
		Protocol: "QUIC",
		Server:   addr,
		Burst:    bunch,
		Start:    time.Now().Format(time.RFC3339),
		Times:    make([]common.Transfer, iter),
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	result.Times[0] = *mkTransfer(stream, buf, 1, iter, time.Now())
	currIter := 1

	if iter > 1 {
		done := make(chan bool, 1)
		defer close(done)
		for {
			select {
			case t := <-ticker.C:
				result.Times[currIter] = *mkTransfer(stream, buf, currIter+1, iter, t)
				currIter++
				if currIter >= iter {
					done <- true
				}
			case <-done:
				common.PrintJSon(result)
				return nil
			}
		}
	}
	return nil
}

func mkTransfer(stream quic.Stream, buf []byte, current int, iter int, t time.Time) *common.Transfer {

	message := fmt.Sprintf("GET %d/%d %d\n", current, iter, len(buf))
	log.Printf("Client: iteration %d, Sending > %s on %v", current, message, stream)

	_, err := stream.Write([]byte(message))
	if common.FatalError(err) != nil {
		return nil
	}

	n, err := io.ReadFull(stream, buf)
	tDelta := time.Since(t).String()

	common.FatalError(err)
	log.Printf("Client: Got %d bytes back\n", n)
	return &common.Transfer{
		XferStart: t.Format(time.RFC3339),
		XferTime:  tDelta,
		XferBytes: n,
		XferIter:  current,
	}
}
