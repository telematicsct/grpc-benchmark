package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/telematicsct/grpc-benchmark/pkg/client"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"
)

func Benchmark_MTLS_HTTP_JSON(b *testing.B) {
	hclient, err := client.NewHTTPSClient()
	if err != nil {
		b.Fatalf("error: %v", err)
	}
	data, err := payload.NewDiagRecorderDataForHTTP()
	if err != nil {
		b.Fatalf("error: %v", err)
	}
	b.ResetTimer()

	// for n := 0; n < b.N; n++ {
	// 	doPost(httpclient, u, b)
	// }

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doPost(hclient, client.GetHttpMTLSUrl(), data, "", b)
		}
	})
}

func doPost(hclient *http.Client, url string, data interface{}, token string, b *testing.B) {
	resp, err := client.DoPost(hclient, url, data, token)
	if err != nil {
		b.Fatalf("http request failed: %v", err)
		return
	}

	defer resp.Body.Close()

	// We need to parse response to have a fair comparison as gRPC does it
	var target payload.DiagResponse
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
