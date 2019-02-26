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

func Benchmark_MTLS_HTTP_HMAC(b *testing.B) {
	client, err := client.NewHTTPSClient()
	if err != nil {
		log.Fatal(err)
	}
	body, err := payload.GetPayload()
	if err != nil {
		b.Fatalf("%v", err)
	}
	u := &mhttp.DiagRecorderData{CanId: payload.GetCanID(), Payload: &mhttp.Payload{Body: body}}
	b.ResetTimer()

	// for n := 0; n < b.N; n++ {
	// 	doHmacPost(httpclient, u, b)
	// }

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doHmacPost(client, u, b)
		}
	})
}

func doHmacPost(hclient *http.Client, data interface{}, b *testing.B) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	req, err := http.NewRequest("POST", client.GetHTTPHmacUrl(), buf)
	if err != nil {
		b.Fatalf("http request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", client.GetJWTToken())

	resp, err := hclient.Do(req)
	if err != nil {
		b.Fatalf("http request failed: %v", err)
		return
	}

	defer resp.Body.Close()

	// We need to parse response to have a fair comparison as gRPC does it
	var target mhttp.DiagResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&target)
	if decodeErr != nil {
		b.Fatalf("unable to decode json: %v", decodeErr)
		return
	}

	if target.Code != 200 || target.Message != "OK" {
		b.Fatalf("http response is wrong: %v", resp)
		return
	}
}
