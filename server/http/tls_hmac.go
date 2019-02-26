package http

import (
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/server"
)

func ServeTLSHMAC(opts *server.ServerOptions) error {
	err := newJWT(opts)
	if err != nil {
		return err
	}
	return doServe(opts, opts.HTTPHMACHostPort, &hmacHandler{}, auth.JWTAuth)
}
