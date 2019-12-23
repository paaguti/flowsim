package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	// "fmt"
	"io/ioutil"
	"log"
	"math/big"
	"path"
)

// Setup a bare-bones TLS config for the server
// Return an empty tls.Config{} on error or empty certs

func ServerTLSConfig(certs string) *tls.Config {
	//
	//  TODO: use flowsimCA.*
	//
	// if certs == "" {
	// 	return &tls.Config{}
	// }
	// 	caCertPEM, err := ioutil.ReadFile(path.Join(certs, "flowsimCA.crt"))
	// 	if FatalErrorln(err, "Reading CA CRT") != nil {
	// 		return &tls.Config{}
	// 	}
	// 	roots := x509.NewCertPool()
	// 	ok := roots.AppendCertsFromPEM(caCertPEM)
	// 	if !ok {
	// 		panic("Failed to parce root certificate")
	// 	}
	// 	cert, err := tls.LoadX509KeyPair(path.Join(certs, "flowsim-server.crt"),
	// 		path.Join(certs, "flowsim-server.key"))
	// 	if FatalError(err) != nil {
	// 		return &tls.Config{}
	// 	}
	// 	return &tls.Config{
	// 		Certificates: []tls.Certificate{cert},
	// 		ClientAuth:   tls.RequireAndVerifyClientCert,
	// 		ClientCAs:    roots,
	// 	}
	//
	// Setup a bare-bones TLS config for the server
	log.Printf("Generating barebones TlSConfig for server, ignoring directory '%s'", certs)
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if FatalError(err) != nil {
		return nil
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if FatalError(err) != nil {
		return nil
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if FatalError(err) != nil {
		return nil
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

// Create a barebones minimum TLS configuration for the client

func ClientTLSConfig(certs string) *tls.Config {
	log.Printf("Ignoring directory '%s' for barebones TLS config", certs)
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

// For HTTPS

func HttpsServerTLSConfig(certs string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(path.Join(certs, "flowsim-client.crt"))
	if err != nil {
		return nil, FatalError(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}, nil
}

func HttpsClientTLSConfig(certs string) (*tls.Config, error) {
	// path.Join(path.Dir(filename), "data.csv")
	caCert, err := ioutil.ReadFile(path.Join(certs, "flowsim-server.crt"))
	if err != nil {
		return nil, FatalError(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair(path.Join(certs, "flowsim-client.crt"), path.Join(certs, "flowsim-client.key"))
	if err != nil {
		return nil, FatalError(err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
