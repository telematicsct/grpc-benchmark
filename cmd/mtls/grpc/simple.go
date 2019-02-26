package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/keepalive"

	"github.com/telematicsct/grpc-benchmark/cmd"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Serve(cliopts *cmd.CliOptions) error {
	dcm := NewDCMServer()
	return goServe(cliopts, cliopts.GRPCHostPort, nil, dcm)
}

// Start starts the grpc server with the provided certificate
func goServe(cliopts *cmd.CliOptions, listen string, option grpc.ServerOption, dcm *dcmServer) error {
	certificate, err := tls.LoadX509KeyPair(cliopts.ServerCertPath, cliopts.ServerKeyPath)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(cliopts.CACertPath)
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

	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true,
		}),
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}
	if option != nil {
		opts = append(opts, option)
	}

	gs := grpc.NewServer(opts...)

	pb.RegisterDCMServiceServer(gs, dcm)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.dcmservice", 1)
	healthpb.RegisterHealthServer(gs, healthServer)

	switch dcm.authType {
	case JWTAuth:
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
