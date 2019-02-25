package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"testing"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
)

func getClient() *grpc.ClientConn {
	certificate, err := tls.LoadX509KeyPair(
		"certs/client.crt",
		"certs/client.key",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		// grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
	}
	conn, err := grpc.Dial("localhost:7900", opts...)
	if err != nil {
		log.Fatalf("failed to dial server: %s", err)
	}

	return conn
}

func getPayload(b *testing.B) []byte {
	//100000 - 100kb
	payload := make([]byte, 100000)
	if _, err := rand.Read(payload); err != nil {
		b.Fatalf("payload error %v", err)
	}
	return payload
}

func Benchmark_MTLS_GRPC_Protobuf(b *testing.B) {
	c := pb.NewDCMServiceClient(getClient())
	payload := getPayload(b)
	data := &pb.DiagRecorderData{CanId: 123456789, Payload: &pb.Payload{Body: payload}}

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
	c := pb.NewDCMServiceClient(getClient())
	payload := getPayload(b)
	data := &pb.DiagRecorderData{CanId: 123456789, Payload: &pb.Payload{Body: payload}}

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
