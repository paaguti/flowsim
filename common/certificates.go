package common

import (
	// "crypto/rand"
	// "crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	// "encoding/pem"
	// "fmt"
	"io/ioutil"
	"path"
	// "math/big"
)

// Setup a bare-bones TLS config for the server
// Return an empty tls.Config{} on error or empty certs

func ServerTLSConfig(certs string) *tls.Config {
	if certs == "" {
		return &tls.Config{}
	}
	// caCert, err := ioutil.ReadFile(path.Join(certs, "client.crt"))
	caCertPEM, err := ioutil.ReadFile(path.Join(certs, "flowsimCA.crt"))
	if FatalErrorln(err, "Reading CA CRT") != nil {
		return &tls.Config{}
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("Failed to parce root certificate")
	}
	cert, err := tls.LoadX509KeyPair(path.Join(certs, "flowsim-server.crt"),
		path.Join(certs, "flowsim-server.key"))
	if FatalError(err) != nil {
		return &tls.Config{}
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}
}

// Create a barebones minimum TLS configuration for the client

func ClientTLSConfig() *tls.Config {

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
