#!/bin/bash

install_dir=$(which flowsim)
if [ -z "$install_dir" ]; then
	echo "Install flowsim FIRST"
	exit
fi
cd $(dirname $install_dir)

openssl req \
    -newkey rsa:2048 \
    -nodes \
    -days 3650 \
    -x509 \
    -keyout flowsimCA.key \
    -out flowsimCA.crt \
    -subj "/CN=*"
openssl req \
    -newkey rsa:2048 \
    -nodes \
    -keyout flowsim-server.key \
    -out flowsim-server.csr \
    -subj "/C=GB/ST=London/L=London/O=Global Security/OU=IT Department/CN=*"
openssl x509 \
    -req \
    -days 365 \
    -sha256 \
    -in flowsim-server.csr \
    -CA flowsimCA.crt \
    -CAkey flowsimCA.key \
    -CAcreateserial \
    -out flowsim-server.crt \
    -extfile <(echo subjectAltName = IP:127.0.0.1)

openssl req \
    -x509 \
    -nodes \
    -newkey rsa:2048 \
    -keyout flowsim-client.key \
    -out flowsim-client.crt \
    -days 3650 \
    -subj "/C=GB/ST=London/L=London/O=Global Security/OU=IT Department/CN=*"
cd -

