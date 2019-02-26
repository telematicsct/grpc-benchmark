package mgrpc

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/telematicsct/grpc-benchmark/cmd"
	"google.golang.org/grpc"
)

func ServeHMAC(cliopts *cmd.CliOptions) error {
	dcm, err := NewDCMServerWithJWT(cliopts.JWTPrivateKey, cliopts.JWTPublicKey)
	if err != nil {
		return err
	}
	opt := grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(jwtAuthFunc))
	return goServe(cliopts, cliopts.GRPCHMACHostPort, opt, dcm)
}
