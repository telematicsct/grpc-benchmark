package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	//APIKey for server
	APIKey = "2730943b-ba0c-4038-bb5f-21965bf24b6d"
	// AuthHeader defines authorization header.
	AuthHeader = "Authorization"
	// AuthScheme defines authorization scheme.
	AuthScheme = "Bearer"
	// AuthorizationKey is the key used to store authorization token data
	AuthorizationKey = "authorization"
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

	myAuthFunction := func(ctx context.Context) (context.Context, error) {

		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
		}
		if err := validate(meta, APIKey); err != nil {
			return nil, err
		}
		// val, err := grpc_auth.AuthFromMD(ctx, AuthHeader)
		// if err != nil {
		// 	return nil, err
		// }
		// if val != APIKey {
		// 	log.Fatalln("invalid api key")
		// 	return nil, grpc.Errorf(codes.Unauthenticated, "Invalid API key")
		// }
		return ctx, nil
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true,
		}),
		grpc.Creds(creds),
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(myAuthFunction)),
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

func validate(meta metadata.MD, key string) error {
	keys, ok := meta[AuthorizationKey]
	if !ok || len(meta[AuthorizationKey]) == 0 {
		return grpc.Errorf(codes.Unauthenticated, "no key provided")
	}
	if keys[0] != key {
		return grpc.Errorf(codes.Unauthenticated, "invalid key")
	}
	return nil
}
