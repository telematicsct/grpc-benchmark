package server

import (
	"github.com/urfave/cli"
)

type ServerOptions struct {
	HTTPHostPort     string
	HTTPHMACHostPort string
	GRPCHostPort     string
	GRPCHMACHostPort string
	ServerCertPath   string
	ServerKeyPath    string
	CACertPath       string
	AuthKey          string
	JWTPrivateKey    string
	JWTPublicKey     string
}

func NewServerOptions(c *cli.Context) *ServerOptions {
	opts := &ServerOptions{}
	opts.HTTPHostPort = c.String("http-listen")
	opts.HTTPHMACHostPort = c.String("http-hmac-listen")
	opts.GRPCHostPort = c.String("grpc-listen")
	opts.GRPCHMACHostPort = c.String("grpc-hmac-listen")
	opts.CACertPath = c.String("ca")
	opts.ServerCertPath = c.String("cert")
	opts.ServerKeyPath = c.String("key")
	opts.AuthKey = c.String("auth-key")
	opts.JWTPrivateKey = c.String("jwt-private-key")
	opts.JWTPublicKey = c.String("jwt-public-key")
	return opts
}
