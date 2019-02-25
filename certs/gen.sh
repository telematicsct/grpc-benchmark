#!/bin/bash

scriptdir=`dirname $0`
cd $scriptdir

cmd=$1
shift

ca_config=ca.conf
config=san.conf
options="-CAcreateserial"
[[ -f ca.srl ]] && options="-CAserial ca.srl"

case $cmd in
  ca)
    [[ -f ca.key && $UPDATE != "y" ]] && echo "CA key already generated." && exit
    # ca.key
    openssl genrsa -out ca.key 2048

    # ca.crt
    openssl req -new -key ca.key -x509 -days 3650 -out ca.crt -subj /C=US/ST=Texas/L=Plano/O=Toyota/CN="Benchmark Root CA" -extensions v3_ca -config $ca_config
    ;;

  server)
    [[ -f server.key && $UPDATE != "y"  ]] && echo "server key already generated." && exit
    # server.key
    openssl genrsa -out server.key 2048

    # server.csr
    openssl req -new -sha256 -key server.key -out server.csr -config $config -extensions v3_req

    # server.crt
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -out server.crt $options -extensions v3_req -extfile $config
    ;;

  server-update)
    # server.csr
    openssl req -new -sha256 -key server.key -out server.csr -config $config -extensions v3_req

    # server.crt
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -out server.crt $options -extensions v3_req -extfile $config
    ;;

  client)
    [[ -z $1 ]] && echo "missing client name arg" && exit
    [[ -f $1.key && $UPDATE != "y" ]] && echo "$1 client key already generated." && exit
    # client.key
    openssl genrsa -out $1.key 2048

    # client.csr
    openssl req -new -nodes -key $1.key -out $1.csr -subj /C=US/ST=Texas/L=Plano/O=Toyota/CN=localhost

    # client.crt
    openssl x509 -req -in $1.csr -CA ca.crt -CAkey ca.key -out $1.crt $options -extfile
    ;;
esac
