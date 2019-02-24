package benchmarks

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"testing"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getClient() *grpc.ClientConn {
	certificate, err := tls.LoadX509KeyPair(
		"certs/client1.crt",
		"certs/client1.key",
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

	dialOption := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial("localhost:7900", dialOption)
	if err != nil {
		log.Fatalf("failed to dial server: %s", err)
	}

	return conn
}

func Benchmark_MTLS_GRPC_Protobuf(b *testing.B) {
	c := pb.NewDCMServiceClient(getClient())
	for n := 0; n < b.N; n++ {
		doGRPC(c, b)
	}
}

func doGRPC(c pb.DCMServiceClient, b *testing.B) {
	stream, err := c.DiagnosticData(context.Background())
	if err != nil {
		b.Fatalf("%v.DiagnosticData(_) = _, %v", c, err)
	}
	payload := make([]byte, 1)
	/*
		_, err = rand.Read(payload)
		if err != nil {
			b.Fatalf("payload error %v", err)
		}
	*/
	data := &pb.DiagRecorderData{CanId: 123456789, Payload: &pb.Payload{Body: payload}}
	if err := stream.Send(data); err != nil {
		b.Fatalf("send error %v", err)
	}

	if err != nil {
		b.Fatalf("grpc request failed: %v", err)
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	if reply.Code != 200 || reply.Message != "Done" {
		b.Fatalf("grpc response is wrong: %v", reply)
	}
}
