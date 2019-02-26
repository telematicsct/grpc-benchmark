package server

import (
	"github.com/urfave/cli"
)

type ServerOptions struct {
	HTTPMTLSHostPort     string
	HTTPTLSHmacHostPort  string
	HTTPMTLSHmacHostPort string
	GRPCHostPort         string
	GRPCHMACHostPort     string
	ServerCertPath       string
	ServerKeyPath        string
	CACertPath           string
	AuthKey              string
	JWTPrivateKey        string
	JWTPublicKey         string
}

func NewServerOptions(c *cli.Context) *ServerOptions {
	opts := &ServerOptions{}
	opts.HTTPMTLSHostPort = c.String("http-mtls-listen")
	opts.HTTPTLSHmacHostPort = c.String("http-tls-hmac-listen")
	opts.HTTPMTLSHmacHostPort = c.String("http-mtls-hmac-listen")
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
