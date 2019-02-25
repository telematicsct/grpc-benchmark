package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/telematicsct/grpc-benchmark/pkg/env"
	"io/ioutil"
	"net/http"
)

const (
	//CertBasePath base path for certificates
	CertBasePath = "CERT_BASE_PATH"
)

var (
	caCertPath     = env.GetString(CertBasePath, "certs/ca.crt")
	clientCertPath = env.GetString(CertBasePath, "certs/client.crt")
	clientKeyPath  = env.GetString(CertBasePath, "certs/client.key")
)

//NewHTTPSClient returns a new https client
func NewHTTPSClient() (*http.Client, error) {

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
