package main

import (
	"log"
	"testing"

	"github.com/telematicsct/grpc-benchmark/pkg/client"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"
)

func Benchmark_MTLS_HTTP_HMAC_JSON(b *testing.B) {
	hclient, err := client.NewHTTPSClient()
	if err != nil {
		log.Fatal(err)
	}
	data, err := payload.NewDiagRecorderDataForHTTP()
	if err != nil {
		b.Fatalf("error: %v", err)
	}
	b.ResetTimer()

	// for n := 0; n < b.N; n++ {
	// 	doHmacPost(httpclient, u, b)
	// }

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doPost(hclient, client.GetHttpMTLSHmacUrl(), data, client.GetJWTToken(), b)
		}
	})
}
