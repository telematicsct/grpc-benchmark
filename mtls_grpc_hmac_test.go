package main

import (
	"log"
	"testing"

	"github.com/telematicsct/grpc-benchmark/pkg/client"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"golang.org/x/net/context"
	_ "google.golang.org/grpc/encoding/gzip"
)

func Benchmark_MTLS_GRPC_HMAC_Protobuf(b *testing.B) {
	c, err := client.NewDCMServiceClientWithAuth(client.GetGRPCHmacUrl(), client.GetJWTToken())
	if err != nil {
		b.Fatalf("%v", err)
	}
	data, err := payload.NewDiagRecorderData()
	if err != nil {
		b.Fatalf("%v", err)
	}
	//warm up
	for i := 0; i < 5; i++ {
		doGRPCHmac(c, data, b)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doGRPCHmac(c, data, b)
		}
	})

	// for n := 0; n < b.N; n++ {
	// 	doGRPCHmac(c, data, b)
	// }
}

func Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream(b *testing.B) {
	c, err := client.NewDCMServiceClientWithAuth(client.GetGRPCHmacUrl(), client.GetJWTToken())
	if err != nil {
		b.Fatalf("%v", err)
	}
	data, err := payload.NewDiagRecorderData()
	if err != nil {
		b.Fatalf("%v", err)
	}
	//warm up
	for i := 0; i < 5; i++ {
		doGRPCHmacStream(c, data, b)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doGRPCHmacStream(c, data, b)
		}
	})

	// for n := 0; n < b.N; n++ {
	// 	doGRPCHmacStream(c, data, b)
	// }
}

func doGRPCHmac(c pb.DCMServiceClient, data *pb.DiagRecorderData, b *testing.B) {
	resp, err := c.DiagnosticData(context.Background(), data)
	if err != nil {
		b.Fatalf("grpc request failed: %v", err)
	}

	if resp == nil || resp.Code != 200 {
		b.Fatalf("wrong grpc response %v", resp)
	}
}

func doGRPCHmacStream(c pb.DCMServiceClient, data *pb.DiagRecorderData, b *testing.B) {
	stream, err := c.DiagnosticDataStream(context.Background())
	if err != nil {
		b.Fatalf("%v.DiagnosticData(_) = _, %v", c, err)
	}

	if err := stream.Send(data); err != nil {
		b.Fatalf("send error %v", err)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	if reply.Code != 200 || reply.Message != "Done" {
		b.Fatalf("grpc response is wrong: %v", reply)
	}
}
