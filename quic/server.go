package quic

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"

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

	//
	// TODO:
	//   Include certificate directory handling ("/etc") is just a dummy
	//
	tlsConfig, err := common.ServerTLSConfig("/etc")
	if common.FatalError(err) != nil {
		return err
	}
	tlsConfig.NextProtos = []string{"flowsim-quic"}
	listener, err := quic.Listen(conn, tlsConfig, nil)
	if common.FatalError(err) != nil {
		return err
	}

	for {
		// This is for the latest version of quic-go
		// stream, err := session.OpenStreamSync(context.Background())
		//
		// revert to get the spin bit running
		//
		sess, err := listener.Accept(context.Background())

		if common.FatalError(err) != nil {
			return err
		}
		if single {
			quicHandler(sess)
			return nil
		}
		go quicHandler(sess)
	}
}

func quicHandler(sess quic.Session) error {

	log.Println("Entering quicHandler")
	//
	// This is for the latest version of quic-go
	stream, err := sess.OpenStreamSync(context.Background())

	if common.FatalError(err) != nil {
		return err
	}
	log.Println("Got a stream")

	// reader := bufio.NewReader(stream)
	cmd := make([]byte, 128)
	for {
		// log.Println("In server loop")
		// cmd, err := reader.ReadString('\n')
		_, err := io.ReadFull(stream, cmd)
		if common.FatalError(err) != nil {
			return err
		}
		log.Printf("In server loop: got %s", cmd)
		wbuf, end, err := parseCmd(string(cmd))
		if common.FatalError(err) != nil {
			return err
		}
		_, err = stream.Write(wbuf)
		if common.FatalError(err) != nil {
			return err
		}
		if end {
			break
		}
	}
	return err
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
	log.Printf("Server: Got %s", strb)
	iter, total, bunchStr, err := matcher(strb)
	if err == nil {
		bunch, _ := strconv.Atoi(bunchStr) // ignore error, wouldn't have parsed the command
		nb := make([]byte, bunch, bunch)
		_, err := rand.Read(nb)
		if err != nil {
			log.Printf("ERROR while filling random buffer: %v\n", err)
			return nil, iter == total, err
		}
		log.Printf("Sending %d bytes\n", len(nb))
		return nb, iter == total, err
	}
	return nil, false, err
}

// Setup a bare-bones TLS config for the server (moded to)
// common.ServerTLSConfig(certs string) *tls.Config
//
// func generateTLSConfig() *tls.Config {
// 	key, err := rsa.GenerateKey(rand.Reader, 1024)
// 	if common.FatalError(err) != nil {
// 		return nil
// 	}
// 	template := x509.Certificate{SerialNumber: big.NewInt(1)}
// 	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
// 	if common.FatalError(err) != nil {
// 		return nil
// 	}
// 	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
// 	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

// 	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
// 	if common.FatalError(err) != nil {
// 		return nil
// 	}
// 	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
// }
