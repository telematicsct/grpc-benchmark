package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/telematicsct/grpc-benchmark/pkg/service"
	"io/ioutil"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/server"
)

//ServeMTLS creates and serves gRPC MTLS server
func ServeMTLS(opts *server.ServerOptions) error {
	dcm := service.NewDCMService()
	return goServe(opts, opts.GRPCHostPort, nil, dcm)
}

// Start starts the grpc server with the provided certificate
func goServe(opts *server.ServerOptions, listen string, grpcoption grpc.ServerOption, dcm *service.DCM) error {
	certificate, err := tls.LoadX509KeyPair(opts.ServerCertPath, opts.ServerKeyPath)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(opts.CACertPath)
	if err != nil {
		return err
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return errors.New("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	grpcopts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true,
		}),
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}
	if grpcoption != nil {
		grpcopts = append(grpcopts, grpcoption)
	}

	gs := grpc.NewServer(grpcopts...)

	pb.RegisterDCMServiceServer(gs, dcm)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.dcmservice", 1)
	healthpb.RegisterHealthServer(gs, healthServer)

	switch dcm.AuthType {
	case auth.JWTAuth:
		log.Println("GRPC MTLS HMAC(JWT) Listening at", listen)
	default:
		log.Println("GRPC MTLS Listening at", listen)
	}

	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	return gs.Serve(ln)
}
