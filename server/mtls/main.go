package main

import (
	"flag"
	"fmt"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	var (
		debugListenAddr = flag.String("debug-listen-addr", "127.0.0.1:7901", "HTTP listen address.")
		listenAddr      = flag.String("listen-addr", "0.0.0.0:7900", "HTTP listen address.")
		tlsCert         = flag.String("tls-cert", "server-cert.pem", "TLS server certificate.")
		tlsKey          = flag.String("tls-key", "server-key.pem", "TLS server key.")
	)
	flag.Parse()

	log.Println("grpc service starting...")

	creds, err := credentials.NewServerTLSFromFile(*tlsCert, *tlsKey)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to setup tls: %v", err))
	}
	opts := []grpc.ServerOption{
		// grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		// 	MinTime:             1 * time.Minute,
		// 	PermitWithoutStream: true,
		// }),
		grpc.Creds(creds),
	}

	gs := grpc.NewServer(opts...)

	dcm := NewDCMServer()
	pb.RegisterDCMServiceServer(gs, dcm)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.dcmservice", 1)
	healthpb.RegisterHealthServer(gs, healthServer)

	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	go gs.Serve(ln)

	log.Println("Hello service started successfully.")
	log.Fatal(http.ListenAndServe(*debugListenAddr, nil))

}
