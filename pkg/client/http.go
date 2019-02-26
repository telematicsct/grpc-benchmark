package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/telematicsct/grpc-benchmark/pkg/env"
)

const (
	//CertBasePath base path for certificates
	CertBasePath = "CERT_BASE_PATH"
	// HTTPURL is http url
	HTTPURL = "HTTP_URL"
)

var (
	httpURL        = env.GetString(HTTPURL, "https://localhost:8443")
	caCertPath     = env.GetString(CertBasePath, "certs/ca.crt")
	clientCertPath = env.GetString(CertBasePath, "certs/client.crt")
	clientKeyPath  = env.GetString(CertBasePath, "certs/client.key")
)

// GetHTTPUrl returns http url
func GetHTTPUrl() string {
	return httpURL
}

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
