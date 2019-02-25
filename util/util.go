package util

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getCertsBasePath() string {
	path := os.Getenv("CERT_BASE_PATH")
	if len(path) == 0 {
		return "certs"
	}
	return path
}

func GetHTTPSClient() (*http.Client, error) {
	caCertPath := filepath.Join(getCertsBasePath(), "ca.crt")
	clientCertPath := filepath.Join(getCertsBasePath(), "client.crt")
	clientKeyPath := filepath.Join(getCertsBasePath(), "client.key")

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
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

func GetGRPCClient(listenAddr string) (*grpc.ClientConn, error) {
	host := listenAddr
	if strings.Contains(listenAddr, ":") {
		parts := strings.Split(listenAddr, ":")
		host = parts[0]
	}
	caCertPath := filepath.Join(getCertsBasePath(), "ca.crt")
	clientCertPath := filepath.Join(getCertsBasePath(), "client.crt")
	clientKeyPath := filepath.Join(getCertsBasePath(), "client.key")

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
		grpc.WithTransportCredentials(transportCreds),
		// grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
	}
	conn, err := grpc.Dial(listenAddr, opts...)
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

func NewDCMServiceClient(listenAddr string) (pb.DCMServiceClient, error) {
	conn, err := GetGRPCClient(listenAddr)
	if err != nil {
		return nil, err
	}
	return pb.NewDCMServiceClient(conn), nil
}
