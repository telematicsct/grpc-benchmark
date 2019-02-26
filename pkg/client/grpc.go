package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/telematicsct/grpc-benchmark/pkg/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	GRPC_URL = "GRPC_URL"
)

var (
	grpcURL = env.GetString(GRPC_URL, "localhost:7900")
)

// GetGRPCUrl grpc url
func GetGRPCUrl() string {
	return grpcURL
}

//NewGRPCClient returns a new gRPC client
func NewGRPCClient(listenAddr string) (*grpc.ClientConn, error) {
	host := listenAddr
	if strings.Contains(listenAddr, ":") {
		parts := strings.Split(listenAddr, ":")
		host = parts[0]
	}

	certificate, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, errors.New("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   host,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(transportCreds),
		// grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
	}
	conn, err := grpc.Dial(listenAddr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
