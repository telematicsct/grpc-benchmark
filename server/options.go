package server

import (
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/urfave/cli"
)

type ProtocolType int

const (
	HTTP ProtocolType = iota
	GRPC
)

type TLSType int

const (
	TLS TLSType = iota
	MTLS
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

// GetBind returns the bind (hostport) string for the protocol, tls, auth type.
func (o *ServerOptions) GetBind(protocol ProtocolType, tls TLSType, authType auth.AuthType) string {
	switch protocol {
	case HTTP:
		switch tls {
		case TLS:
			switch authType {
			case auth.JWTAuth:
				return o.HTTPTLSHmacHostPort
			}
		case MTLS:
			switch authType {
			case auth.NoAuth:
				return o.HTTPMTLSHostPort
			case auth.JWTAuth:
				return o.HTTPMTLSHmacHostPort
			}
		}
	case GRPC:
		switch tls {
		case MTLS:
			switch authType {
			case auth.NoAuth:
				return o.GRPCHostPort
			case auth.JWTAuth:
				return o.GRPCHMACHostPort
			}
		}
	}

	return ""
}
