package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/telematicsct/grpc-benchmark/server"
	grpcServer "github.com/telematicsct/grpc-benchmark/server/grpc"
	httpServer "github.com/telematicsct/grpc-benchmark/server/http"
)

func main() {
	app := cli.NewApp()
	httpMTLSSimpleListenFlag := cli.StringFlag{
		Name:  "http-mtls-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:8443",
	}
	httpMTLSHmacListenFlag := cli.StringFlag{
		Name:  "http-mtls-hmac-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:9443",
	}
	httpTLSHmacListenFlag := cli.StringFlag{
		Name:  "http-tls-hmac-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:7443",
	}
	grpcListenFlag := cli.StringFlag{
		Name:  "grpc-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:7900",
	}
	grpcHmacListenFlag := cli.StringFlag{
		Name:  "grpc-hmac-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:8900",
	}
	certFlag := cli.StringFlag{
		Name:  "cert, c",
		Usage: "tls certificate",
		Value: "certs/server.crt",
	}
	keyFlag := cli.StringFlag{
		Name:  "key, k",
		Usage: "tls key",
		Value: "certs/server.key",
	}
	caFlag := cli.StringFlag{
		Name:  "ca",
		Usage: "ca key",
		Value: "certs/ca.crt",
	}
	jwtPublicKeyFlag := cli.StringFlag{
		Name:  "jwt-public-key",
		Usage: "jwt public key",
		Value: "certs/jwt.pub.pem",
	}
	jwtPrivateKeyFlag := cli.StringFlag{
		Name:  "jwt-private-key",
		Usage: "jwt private key",
		Value: "certs/jwt",
	}
	app.Commands = []cli.Command{
		{
			Name:  "all",
			Usage: "all",
			Flags: []cli.Flag{
				httpMTLSSimpleListenFlag, httpMTLSHmacListenFlag, httpTLSHmacListenFlag,
				grpcListenFlag, grpcHmacListenFlag,
				jwtPrivateKeyFlag, jwtPublicKeyFlag,
				certFlag, keyFlag, caFlag,
			},
			Action: func(c *cli.Context) error {
				opts := server.NewServerOptions(c)
				go func() {
					if err := httpServer.ServeMTLS(opts); err != nil {
						log.Fatalf("failed to start http mtls server: %s", err)
					}
				}()

				go func() {
					if err := httpServer.ServeMTLSHMAC(opts); err != nil {
						log.Fatalf("failed to start http mtls (HMAC) server: %s", err)
					}
				}()

				go func() {
					if err := httpServer.ServeTLSHMAC(opts); err != nil {
						log.Fatalf("failed to start http tls (HMAC) server: %s", err)
					}
				}()

				go func() {
					if err := grpcServer.ServeMTLS(opts); err != nil {
						log.Fatalf("failed to start gRPC mtls server: %s", err)
					}
				}()

				go func() {
					if err := grpcServer.ServeMTLSHMAC(opts); err != nil {
						log.Fatalf("failed to start gRPC mtls (HMAC) server: %s", err)
					}
				}()

				select {}
			},
		},
		{
			Name:  "http-mtls",
			Usage: "http-mtls",
			Flags: []cli.Flag{httpMTLSSimpleListenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return httpServer.ServeMTLS(server.NewServerOptions(c))
			},
		},
		{
			Name:  "http-mtls-hmac",
			Usage: "http-mtls-hmac",
			Flags: []cli.Flag{httpMTLSHmacListenFlag, jwtPrivateKeyFlag, jwtPublicKeyFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return httpServer.ServeMTLS(server.NewServerOptions(c))
			},
		},
		{
			Name:  "http-tls-hmac",
			Usage: "http-tls-hmac",
			Flags: []cli.Flag{httpTLSHmacListenFlag, jwtPrivateKeyFlag, jwtPublicKeyFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return httpServer.ServeMTLS(server.NewServerOptions(c))
			},
		},
		{
			Name:  "grpc",
			Usage: "grpc",
			Flags: []cli.Flag{grpcListenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return grpcServer.ServeMTLS(server.NewServerOptions(c))
			},
		},
		{
			Name:  "grpc-hmac",
			Usage: "grpc-hmac",
			Flags: []cli.Flag{grpcHmacListenFlag, jwtPrivateKeyFlag, jwtPublicKeyFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return grpcServer.ServeMTLSHMAC(server.NewServerOptions(c))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
