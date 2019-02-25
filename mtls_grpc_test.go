package main

import (
	"github.com/telematicsct/grpc-benchmark/pkg/client"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"
	"log"
	"testing"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"golang.org/x/net/context"
	_ "google.golang.org/grpc/encoding/gzip"
)

func Benchmark_MTLS_GRPC_Protobuf(b *testing.B) {
	c, err := client.NewDCMServiceClient("a3bae774238fe11e9b4530aa49b34ad2-baa821165bd29c97.elb.ap-northeast-1.amazonaws.com:7900")
	if err != nil {
		b.Fatalf("%v", err)
	}
	data, err := payload.NewDiagRecorderData()
	if err != nil {
		b.Fatalf("%v", err)
	}
	//warm up
	for i := 0; i < 5; i++ {
		doGRPC(c, data, b)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doGRPC(c, data, b)
		}
	})

	// for n := 0; n < b.N; n++ {
	// 	doGRPC(c, data, b)
	// }
}

func Benchmark_MTLS_GRPC_Protobuf_Stream(b *testing.B) {
	c, err := client.NewDCMServiceClient("a3bae774238fe11e9b4530aa49b34ad2-baa821165bd29c97.elb.ap-northeast-1.amazonaws.com:7900")
	if err != nil {
		b.Fatalf("%v", err)
	}
	data, err := payload.NewDiagRecorderData()
	if err != nil {
		b.Fatalf("%v", err)
	}
	//warm up
	for i := 0; i < 5; i++ {
		doGRPCStream(c, data, b)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doGRPCStream(c, data, b)
		}
	})

	// for n := 0; n < b.N; n++ {
	// 	doGRPCStream(c, data, b)
	// }
}

func doGRPC(c pb.DCMServiceClient, data *pb.DiagRecorderData, b *testing.B) {
	resp, err := c.DiagnosticData(context.Background(), data)
	if err != nil {
		b.Fatalf("grpc request failed: %v", err)
	}

	if resp == nil || resp.Code != 200 {
		b.Fatalf("wrong grpc response %v", resp)
	}
}

func doGRPCStream(c pb.DCMServiceClient, data *pb.DiagRecorderData, b *testing.B) {
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
