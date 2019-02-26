package main

import (
	"log"
	"os"

	"github.com/telematicsct/grpc-benchmark/cmd"
	"github.com/telematicsct/grpc-benchmark/cmd/mgrpc"
	"github.com/telematicsct/grpc-benchmark/cmd/mhttp"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	httpsListenFlag := cli.StringFlag{
		Name:  "http-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:8443",
	}
	httpsHmacListenFlag := cli.StringFlag{
		Name:  "http-hmac-listen",
		Usage: "Listen address",
		Value: "0.0.0.0:8553",
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
				httpsListenFlag, httpsHmacListenFlag,
				grpcListenFlag, grpcHmacListenFlag,
				jwtPrivateKeyFlag, jwtPublicKeyFlag,
				certFlag, keyFlag, caFlag,
			},
			Action: func(c *cli.Context) error {
				cliopts := cmd.NewCliOptions(c)
				go func() {
					if err := mhttp.Serve(cliopts); err != nil {
						log.Fatalf("failed to start http mtls server: %s", err)
					}
				}()

				go func() {
					if err := mgrpc.Serve(cliopts); err != nil {
						log.Fatalf("failed to start gRPC mtls server: %s", err)
					}
				}()

				go func() {
					if err := mgrpc.ServeHMAC(cliopts); err != nil {
						log.Fatalf("failed to start gRPC mtls server: %s", err)
					}
				}()

				select {}
			},
		},
		{
			Name:  "https",
			Usage: "https",
			Flags: []cli.Flag{httpsListenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return mhttp.Serve(cmd.NewCliOptions(c))
			},
		},
		{
			Name:  "https-hmac",
			Usage: "https-hmac",
			Flags: []cli.Flag{httpsHmacListenFlag, jwtPrivateKeyFlag, jwtPublicKeyFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return mhttp.Serve(cmd.NewCliOptions(c))
			},
		},
		{
			Name:  "grpc",
			Usage: "grpc",
			Flags: []cli.Flag{grpcListenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return mgrpc.Serve(cmd.NewCliOptions(c))
			},
		},
		{
			Name:  "grpc-hmac",
			Usage: "grpc-hmac",
			Flags: []cli.Flag{grpcHmacListenFlag, jwtPrivateKeyFlag, jwtPublicKeyFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return mgrpc.ServeHMAC(cmd.NewCliOptions(c))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
