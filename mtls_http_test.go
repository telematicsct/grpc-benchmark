package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	mtlshttp "github.com/telematicsct/grpc-benchmark/cmd/https"
	"github.com/telematicsct/grpc-benchmark/util"
)

var httpclient *http.Client

func init() {
	client, err := util.GetHTTPSClient()
	if err != nil {
		log.Fatal(err)
	}
	httpclient = client
}

func Benchmark_MTLS_HTTP(b *testing.B) {
	body, err := util.GetPayload()
	if err != nil {
		b.Fatalf("%v", err)
	}
	u := &mtlshttp.DiagRecorderData{CanId: util.GetCanId(), Payload: &mtlshttp.Payload{Body: body}}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		doPost(httpclient, u, b)
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
