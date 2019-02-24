package main

import (
	"context"
	"crypto/rand"
	"flag"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

func main() {

	var (
		serverAddr = flag.String("server-addr", "localhost:7900", "Hello service address.")
		tlsCert    = flag.String("tls-cert", "../certs/server-cert.pem", "TLS server certificate.")
		// token      = flag.String("token", ".token", "Path to Hmac/JWT auth token.")
	)
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile(*tlsCert, "")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(*serverAddr,
		grpc.WithTransportCredentials(creds),
		// grpc.WithPerRPCCredentials(jwtCreds),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewDCMServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := c.DiagnosticData(ctx)
	if err != nil {
		log.Fatalf("%v.DiagnosticData(_) = _, %v", c, err)
	}

	payload := make([]byte, 500000)
	_, err = rand.Read(payload)
	if err != nil {
		log.Fatalln("payload error", err)
	}

	//streaming Diag Recorder Data
	data := &pb.DiagRecorderData{
		CanId:   123456789,
		Payload: &pb.Payload{Body: payload},
	}
	if err := stream.Send(data); err != nil {
		log.Fatalf("%v.Send(%v) = %v", stream, data, err)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("DiagnosticData summary: %v => %v", reply.Code, reply.Message)

}
