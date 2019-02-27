package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"time"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/pkg/service"
	"github.com/telematicsct/grpc-benchmark/server"
)

func Serve(opts *server.ServerOptions, tlsType server.TLSType, authType auth.AuthType) error {
	listen := opts.GetBind(server.GRPC, tlsType, authType)

	var dcm *service.DCM
	var err error

	switch authType {
	case auth.JWTAuth:
		dcm, err = service.NewDCMServiceWithJWT(opts.JWTPrivateKey, opts.JWTPublicKey)
		if err != nil {
			return err
		}
	default:
		dcm = service.NewDCMService()
	}

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

	switch authType {
	case auth.JWTAuth:
		grpcopts = append(grpcopts, grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.JWTAuthFunc(dcm.Token))))
		log.Println("GRPC MTLS HMAC(JWT) Listening at", listen)
	default:
		log.Println("GRPC MTLS Listening at", listen)
	}

	gs := grpc.NewServer(grpcopts...)

	pb.RegisterDCMServiceServer(gs, dcm)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.dcmservice", 1)
	grpc_health_v1.RegisterHealthServer(gs, healthServer)

	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	return gs.Serve(ln)
}
