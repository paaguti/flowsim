package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	// "fmt"
	// "io/ioutil"
	"log"
	"math/big"
	// "path"
)

// Setup a bare-bones TLS config for the server
// Return an empty tls.Config{} on error or empty certs

func ServerTLSConfig(certs string, nextProto string) (*tls.Config, error) {
	// Setup a bare-bones TLS config for the server
	log.Printf("Generating barebones TlSConfig for server, ignoring directory '%s'", certs)
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if FatalError(err) != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if FatalError(err) != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if FatalError(err) != nil {
		return nil, err
	}
	if nextProto == "" {
		return &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		}, nil
	} else {
		return &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			NextProtos:   []string{nextProto},
		}, nil
	}
}

// Create a barebones minimum TLS configuration for the client

func ClientTLSConfig(certs string, nextProto string) (*tls.Config, error) {
	log.Printf("Ignoring directory '%s' for barebones TLS config", certs)
	// log.Println("H2QUIC client test...")
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if nextProto == "" {
		return &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
		}, nil
	} else {
		return &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			NextProtos:         []string{nextProto},
		}, nil
	}
}

// For HTTPS

func HttpsServerTLSConfig(certs string) (*tls.Config, error) {
	log.Printf("HttpsServerTLSConfig(%s)", certs)

	return &tls.Config{
		ClientAuth: tls.RequestClientCert,
	}, nil
}

func HttpsClientTLSConfig(certs string) (*tls.Config, error) {
	log.Printf("HttpsClientTLSConfig(%s)", certs)

	return &tls.Config{InsecureSkipVerify: true}, nil
}

func IsSecureConfig(tlsConfig *tls.Config) bool {

	if tlsConfig.InsecureSkipVerify == true {
		return true
	}
	if tlsConfig.Certificates != nil {
		return true
	}
	if tlsConfig.ClientAuth == tls.RequestClientCert {
		return true
	}
	return false
}
