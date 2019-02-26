package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/telematicsct/grpc-benchmark/cmd/mhttp"
	"github.com/telematicsct/grpc-benchmark/pkg/client"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"
)

var httpclient *http.Client

func init() {
	client, err := client.NewHTTPSClient()
	if err != nil {
		log.Fatal(err)
	}
	httpclient = client
}

func Benchmark_MTLS_HTTP(b *testing.B) {
	body, err := payload.GetPayload()
	if err != nil {
		b.Fatalf("%v", err)
	}
	u := &mhttp.DiagRecorderData{CanId: payload.GetCanID(), Payload: &mhttp.Payload{Body: body}}
	b.ResetTimer()

	// for n := 0; n < b.N; n++ {
	// 	doPost(httpclient, u, b)
	// }

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doPost(httpclient, u, b)
		}
	})
}

func doPost(hclient *http.Client, data interface{}, b *testing.B) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	resp, err := hclient.Post(client.GetHTTPUrl(), "application/json", buf)
	if err != nil {
		b.Fatalf("http request failed: %v", err)
	}

	defer resp.Body.Close()

	// We need to parse response to have a fair comparison as gRPC does it
	var target mhttp.DiagResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&target)
	if decodeErr != nil {
		b.Fatalf("unable to decode json: %v", decodeErr)
	}

	if target.Code != 200 || target.Message != "OK" {
		b.Fatalf("http response is wrong: %v", resp)
	}
}
