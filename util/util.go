package util

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetHTTPSClient() (*http.Client, error) {
	caCert, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	return client, nil
}

func GetGRPCClient() (*grpc.ClientConn, error) {
	certificate, err := tls.LoadX509KeyPair(
		"certs/client.crt",
		"certs/client.key",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		return nil, err
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, errors.New("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		// grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
	}
	conn, err := grpc.Dial("localhost:7900", opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetPayload() ([]byte, error) {
	//100000 - 100kb
	payload := make([]byte, 100000)
	if _, err := rand.Read(payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func GetCanId() int32 {
	return 123456789
}

func NewDiagRecorderData() (*pb.DiagRecorderData, error) {
	payload, err := GetPayload()
	if err != nil {
		return nil, err
	}
	data := &pb.DiagRecorderData{CanId: GetCanId(), Payload: &pb.Payload{Body: payload}}
	return data, nil
}

func NewDCMServiceClient() (pb.DCMServiceClient, error) {
	conn, err := GetGRPCClient()
	if err != nil {
		return nil, err
	}
	return pb.NewDCMServiceClient(conn), nil
}
