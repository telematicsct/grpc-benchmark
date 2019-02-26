package cmd

import (
	"github.com/urfave/cli"
)

type CliOptions struct {
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

func NewCliOptions(c *cli.Context) *CliOptions {
	cliopts := &CliOptions{}
	cliopts.HTTPHostPort = c.String("http-listen")
	cliopts.HTTPHMACHostPort = c.String("http-hmac-listen")
	cliopts.GRPCHostPort = c.String("grpc-listen")
	cliopts.GRPCHMACHostPort = c.String("grpc-hmac-listen")
	cliopts.CACertPath = c.String("ca")
	cliopts.ServerCertPath = c.String("cert")
	cliopts.ServerKeyPath = c.String("key")
	cliopts.AuthKey = c.String("auth-key")
	cliopts.JWTPrivateKey = c.String("jwt-private-key")
	cliopts.JWTPublicKey = c.String("jwt-public-key")
	return cliopts
}
