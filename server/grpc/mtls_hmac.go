package grpc

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/pkg/service"
	"github.com/telematicsct/grpc-benchmark/server"
	"google.golang.org/grpc"
)

//ServeMTLSHMAC serves gRPC mtls with HMAC
func ServeMTLSHMAC(opts *server.ServerOptions) error {
	dcm, err := service.NewDCMServiceWithJWT(opts.JWTPrivateKey, opts.JWTPublicKey)
	if err != nil {
		return err
	}
	opt := grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.JWTAuthFunc(dcm.Token)))
	return goServe(opts, opts.GRPCHMACHostPort, opt, dcm)
}
