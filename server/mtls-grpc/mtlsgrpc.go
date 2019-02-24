package mtlsgrpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Start(listen string, cert string, key string, ca string) error {
	log.Println("grpc service starting...")

	certificate, err := tls.LoadX509KeyPair(cert, key)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(ca)
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
		// grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		// 	MinTime:             1 * time.Minute,
		// 	PermitWithoutStream: true,
		// }),
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}

	gs := grpc.NewServer(opts...)

	dcm := NewDCMServer()
	pb.RegisterDCMServiceServer(gs, dcm)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.dcmservice", 1)
	healthpb.RegisterHealthServer(gs, healthServer)

	log.Println("Listening at", listen)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	//log.Println("Hello service started successfully.")
	//log.Fatal(http.ListenAndServe(*debugListenAddr, nil))
	return gs.Serve(ln)
}
