package client

import (
	"context"
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
	GRPC_URL      = "GRPC_URL"
	GRPC_HMAC_URL = "GRPC_HMAC_URL"
)

var (
	grpcURL     = env.GetString(GRPC_URL, "localhost:7900")
	grpcHmacURL = env.GetString(GRPC_HMAC_URL, "localhost:8900")
)

// GetGRPCUrl grpc url
func GetGRPCUrl() string {
	return grpcURL
}

func GetGRPCHmacUrl() string {
	return grpcHmacURL
}

//NewGRPCClient returns a new gRPC client
func NewGRPCClient(listenAddr string, token string) (*grpc.ClientConn, error) {
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

	if token != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(NewTokenAuth(token)))
	}

	conn, err := grpc.Dial(listenAddr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type tokenAuth struct {
	token string
}

// New holds per-rpc metadata for the gRPC clients
func NewTokenAuth(token string) credentials.PerRPCCredentials {
	return tokenAuth{token}
}

func (j tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j tokenAuth) RequireTransportSecurity() bool {
	return true
}
