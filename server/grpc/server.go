package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/pkg/service"
	"github.com/telematicsct/grpc-benchmark/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"io/ioutil"
	"log"
	"net"
	"time"
)

//ServeMTLSHMAC serves gRPC mtls with HMAC
func ServeMTLSHMAC(serverOpts *server.ServerOptions) error {
	dcm, err := service.NewDCMServiceWithJWT(serverOpts.JWTPrivateKey, serverOpts.JWTPublicKey)
	if err != nil {
		return err
	}
	grpcOpts := grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.JWTAuthFunc(dcm.Token)))
	server := NewDCMServer(serverOpts, grpcOpts, serverOpts.GRPCHMACHostPort, dcm)
	return server.Serve()
}

//ServeMTLS creates and serves gRPC MTLS server
func ServeMTLS(opts *server.ServerOptions) error {
	dcm := service.NewDCMService()
	server := NewDCMServer(opts, nil, opts.GRPCHostPort, dcm)
	return server.Serve()
}

//Server gRPC server
type Server struct {
	ServerOpts *server.ServerOptions
	GRPCOpts   grpc.ServerOption
	HostPort   string
	Service    *service.DCM
}

//NewDCMServer creates and returns DCM server witht he given server and gRPC Options
func NewDCMServer(serverOpts *server.ServerOptions, grpcOpts grpc.ServerOption, hostPort string, dcm *service.DCM) *Server {
	return &Server{serverOpts, grpcOpts, hostPort, dcm}
}

// Serve starts the grpc server with the provided server and gRPC options
func (s *Server) Serve() error {

	certificate, err := tls.LoadX509KeyPair(s.ServerOpts.ServerCertPath, s.ServerOpts.ServerKeyPath)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(s.ServerOpts.CACertPath)
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
	if s.GRPCOpts != nil {
		grpcopts = append(grpcopts, s.GRPCOpts)
	}

	gs := grpc.NewServer(grpcopts...)

	pb.RegisterDCMServiceServer(gs, s.Service)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.dcmservice", 1)
	healthpb.RegisterHealthServer(gs, healthServer)

	switch s.Service.AuthType {
	case auth.JWTAuth:
		log.Println("GRPC MTLS HMAC(JWT) Listening at", s.HostPort)
	default:
		log.Println("GRPC MTLS Listening at", s.HostPort)
	}

	ln, err := net.Listen("tcp", s.HostPort)
	if err != nil {
		return err
	}

	return gs.Serve(ln)
}
