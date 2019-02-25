package main

import (
	"log"
	"testing"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/util"
	"golang.org/x/net/context"
	_ "google.golang.org/grpc/encoding/gzip"
)

func Benchmark_MTLS_GRPC_Protobuf(b *testing.B) {
	c, err := util.NewDCMServiceClient("localhost:7900")
	if err != nil {
		b.Fatalf("%v", err)
	}
	data, err := util.NewDiagRecorderData()
	if err != nil {
		b.Fatalf("%v", err)
	}
	//warm up
	for i := 0; i < 5; i++ {
		doGRPC(c, data, b)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		doGRPC(c, data, b)
	}
}

func Benchmark_MTLS_GRPC_Protobuf_Stream(b *testing.B) {
	c, err := util.NewDCMServiceClient("localhost:7900")
	if err != nil {
		b.Fatalf("%v", err)
	}
	data, err := util.NewDiagRecorderData()
	if err != nil {
		b.Fatalf("%v", err)
	}
	//warm up
	for i := 0; i < 5; i++ {
		doGRPCStream(c, data, b)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		doGRPCStream(c, data, b)
	}
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
