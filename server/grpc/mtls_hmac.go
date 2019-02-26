package grpc

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/telematicsct/grpc-benchmark/server"
	"google.golang.org/grpc"
)

func ServeMTLSHMAC(opts *server.ServerOptions) error {
	dcm, err := NewDCMServerWithJWT(opts.JWTPrivateKey, opts.JWTPublicKey)
	if err != nil {
		return err
	}
	opt := grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(jwtAuthFunc))
	return goServe(opts, opts.GRPCHMACHostPort, opt, dcm)
}
