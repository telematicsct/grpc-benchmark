package main

import (
	"log"
	"os"

	grpc "github.com/telematicsct/grpc-benchmark/cmd/grpc"
	mtlshttp "github.com/telematicsct/grpc-benchmark/cmd/https"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	httpslistenFlag := cli.StringFlag{
		Name:  "httpHostPort",
		Usage: "Listen address",
		Value: "0.0.0.0:8443",
	}
	grpcMTLSListenFlag := cli.StringFlag{
		Name:  "grpcMtlsHostPort",
		Usage: "Listen address",
		Value: "0.0.0.0:7900",
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
	app.Commands = []cli.Command{
		{
			Name:  "all",
			Usage: "all",
			Flags: []cli.Flag{httpslistenFlag, grpcMTLSListenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				go func() {
					if err := mtlshttp.ServerMTLS(c.String("httpHostPort"), c.String("cert"), c.String("key"), c.String("ca")); err != nil {
						log.Fatalf("failed to start http server: %s", err)
					}
				}()

				go func() {
					if err := grpc.ServeMTLS(c.String("grpcMtlsHostPort"), c.String("cert"), c.String("key"), c.String("ca")); err != nil {
						log.Fatalf("failed to start gRPC server: %s", err)
					}
				}()

				select {}
			},
		},
		{
			Name:  "https",
			Usage: "https",
			Flags: []cli.Flag{httpslistenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return mtlshttp.ServerMTLS(c.String("httpHostPort"), c.String("cert"), c.String("key"), c.String("ca"))
			},
		},
		{
			Name:  "grpc",
			Usage: "grpc",
			Flags: []cli.Flag{grpcMTLSListenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return grpc.ServeMTLS(c.String("grpcMtlsHostPort"), c.String("cert"), c.String("key"), c.String("ca"))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
