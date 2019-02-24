#!/bin/bash

scriptdir=`basedir $0`
cd $scriptdir

cmd=$1
shift

case $cmd in
  ca)
    [[ -f ca.key ]] && echo "CA key already generated." && exit
    # ca.key
    openssl genrsa -out ca.key 2048

    # ca.crt
    openssl req -new -key ca.key -x509 -days 3650 -out ca.crt -subj /C=US/ST=Texas/L=Plano/O=Toyota/CN="Benchmark Root CA"
    ;;

  server)
    [[ -f server.key ]] && echo "server key already generated." && exit
    # server.key
    openssl genrsa -out server.key 2048

    # server.csr
    openssl req -new -nodes -key server.key -out server.csr -subj /C=US/ST=Texas/L=Plano/O=Toyota/CN=localhost

    # server.crt
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
    ;;

  client)
    [[ -z $1 ]] && echo "missing client name arg" && exit
    [[ -f $1.key ]] && echo "$1 client key already generated." && exit
    # client.key
    openssl genrsa -out $1.key 2048

    # client.csr
    openssl req -new -nodes -key $1.key -out $1.csr -subj /C=US/ST=Texas/L=Plano/O=Toyota/CN=localhost

    # client.crt
    openssl x509 -req -in $1.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $1.crt
    ;;
esac
