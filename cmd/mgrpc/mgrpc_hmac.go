package mgrpc

import (
	"context"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/telematicsct/grpc-benchmark/cmd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var APIKey = "2730943b-ba0c-4038-bb5f-21965bf24b6d"

const (
	// AuthHeader defines authorization header.
	AuthHeader = "Authorization"
	// AuthScheme defines authorization scheme.
	AuthScheme = "Bearer"
	// AuthorizationKey is the key used to store authorization token data
	AuthorizationKey = "authorization"
)

func ServeHMAC(cliopts *cmd.CliOptions) error {
	APIKey = cliopts.AuthKey
	dcm, err := NewDCMServerWithJWT(cliopts.JWTPrivateKey, cliopts.JWTPublicKey)
	if err != nil {
		return err
	}
	log.Println("xxxListening at", cliopts.GRPCHMACHostPort)
	opt := grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(myAuthFunction))
	return goServe(cliopts, cliopts.GRPCHMACHostPort, opt, dcm)
}

func myAuthFunction(ctx context.Context) (context.Context, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	keys, ok := meta[AuthorizationKey]
	if !ok || len(meta[AuthorizationKey]) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "no key provided")
	}
	if keys[0] != APIKey {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid key")
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
