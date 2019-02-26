package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/telematicsct/grpc-benchmark/pkg/env"
)

const (
	//CertBasePath base path for certificates
	CertBasePath = "CERT_BASE_PATH"
	// HTTPURL is http url
	TARGET_HOST        = "TARGET_HOST"
	HTTP_MTLS_URL      = "HTTP_MTLS_URL"
	HTTP_MTLS_HMAC_URL = "HTTP_MTLS_HMAC_URL"
	HTTP_TLS_HMAC_URL  = "HTTP_TLS_HMAC_URL"
)

var (
	targetHost        = env.GetString(TARGET_HOST, "localhost")
	httpMTLSSimpleURL = env.GetString(HTTP_MTLS_URL, "https://"+TARGET_HOST+":8443")
	httpMTLSHmacURL   = env.GetString(HTTP_MTLS_HMAC_URL, "https://"+TARGET_HOST+":9443")
	httpTLSHmacURL    = env.GetString(HTTP_TLS_HMAC_URL, "https://"+TARGET_HOST+":7443")
	caCertPath        = env.GetString(CertBasePath, "certs/ca.crt")
	clientCertPath    = env.GetString(CertBasePath, "certs/client.crt")
	clientKeyPath     = env.GetString(CertBasePath, "certs/client.key")
	jwtTokenPath      = env.GetString(CertBasePath, "certs/jwt.token")
)

// GetHttpUrl returns http url
func GetHttpMTLSUrl() string {
	return httpMTLSSimpleURL
}

func GetHttpMTLSHmacUrl() string {
	return httpMTLSHmacURL
}

func GetHttpTLSHmacUrl() string {
	return httpTLSHmacURL
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

func GetJWTToken() string {
	data, err := ioutil.ReadFile(jwtTokenPath)
	if err != nil {
		return ""
	}
	return string(data)
}

func DoPost(client *http.Client, url string, data interface{}, token string) (*http.Response, error) {
	if client == nil {
		hclient, err := NewHTTPSClient()
		if err != nil {
			return nil, err
		}
		client = hclient
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	return client.Do(req)
}
