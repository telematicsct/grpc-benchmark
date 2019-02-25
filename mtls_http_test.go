package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	mtlshttp "github.com/telematicsct/grpc-benchmark/cmd/https"
)

var client *http.Client

func init() {
	caCert, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
}

func Benchmark_MTLS_HTTP(b *testing.B) {
	u := &mtlshttp.DiagRecorderData{CanId: 11111, Payload: &mtlshttp.Payload{Body: getPayload(b)}}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		doPost(client, u, b)
	}
}

func doPost(client *http.Client, data interface{}, b *testing.B) {

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	resp, err := client.Post("https://localhost:8443/", "application/json", buf)
	if err != nil {
		b.Fatalf("http request failed: %v", err)
	}

	defer resp.Body.Close()

	// We need to parse response to have a fair comparison as gRPC does it
	var target mtlshttp.DiagResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&target)
	if decodeErr != nil {
		b.Fatalf("unable to decode json: %v", decodeErr)
	}

	if target.Code != 200 || target.Message != "OK" {
		b.Fatalf("http response is wrong: %v", resp)
	}
}
