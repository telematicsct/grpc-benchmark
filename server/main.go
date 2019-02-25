package main

import (
	"log"
	"os"

	grpc "github.com/telematicsct/grpc-benchmark/server/grpc"
	mtlshttp "github.com/telematicsct/grpc-benchmark/server/https"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	httpslistenFlag := cli.StringFlag{
		Name:  "listen, l",
		Usage: "Listen address",
		Value: "0.0.0.0:8443",
	}
	grpclistenFlag := cli.StringFlag{
		Name:  "listen, l",
		Usage: "Listen address",
		Value: "0.0.0.0:7900",
	}
	certFlag := cli.StringFlag{
		Name:  "cert, c",
		Usage: "tls certificate",
		Value: "../certs/server.crt",
	}
	keyFlag := cli.StringFlag{
		Name:  "key, k",
		Usage: "tls key",
		Value: "../certs/server.key",
	}
	caFlag := cli.StringFlag{
		Name:  "ca",
		Usage: "ca key",
		Value: "../certs/ca.crt",
	}
	app.Commands = []cli.Command{
		{
			Name:  "https",
			Usage: "https",
			Flags: []cli.Flag{httpslistenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return mtlshttp.ServerMTLS(c.String("listen"), c.String("cert"), c.String("key"), c.String("ca"))
			},
		},
		{
			Name:  "grpc",
			Usage: "grpc",
			Flags: []cli.Flag{grpclistenFlag, certFlag, keyFlag, caFlag},
			Action: func(c *cli.Context) error {
				return grpc.ServeMTLS(c.String("listen"), c.String("cert"), c.String("key"), c.String("ca"))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
