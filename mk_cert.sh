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
    -keyout ca.key \
    -out ca.crt \
    -subj "/CN=*"
openssl req \
    -newkey rsa:2048 \
    -nodes \
    -keyout server.key \
    -out server.csr \
    -subj "/C=GB/ST=London/L=London/O=Global Security/OU=IT Department/CN=*"
openssl x509 \
    -req \
    -days 365 \
    -sha256 \
    -in server.csr \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out server.crt \
    -extfile <(echo subjectAltName = IP:127.0.0.1)

openssl req \
    -x509 \
    -nodes \
    -newkey rsa:2048 \
    -keyout client.key \
    -out client.crt \
    -days 3650 \
    -subj "/C=GB/ST=London/L=London/O=Global Security/OU=IT Department/CN=*"
cd -

