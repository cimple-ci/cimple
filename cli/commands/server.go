package cli

import (
	"crypto/tls"
	"fmt"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/cimple/server"
	"github.com/urfave/cli"
	"os"
)

func Server() cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Start the Cimple server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "port",
				Usage: "The server port.",
				Value: "8080",
			},
			cli.StringFlag{
				Name:  "host",
				Usage: "The host to bind to",
				Value: "127.0.0.1",
			},
			cli.BoolFlag{
				Name:  "no-tls",
				Usage: "Disable TLS for the server",
			},
			cli.StringFlag{
				Name:  "tls-cert-file",
				Usage: "Specifies the path to the TLS certificate file",
				Value: "server.crt",
			},
			cli.StringFlag{
				Name:  "tls-key-file",
				Usage: "Specifies the path to the TLS key file",
				Value: "server.key",
			},
		},
		Action: func(c *cli.Context) error {
			logging.SetDefaultLogger("Server", os.Stdout)
			logger := logging.CreateLogger("Server", os.Stdout)

			serverConfig := server.DefaultConfig()
			serverConfig.EnableTLS = !c.Bool("no-tls")

			if serverConfig.EnableTLS {
				tlsCertFile := c.String("tls-cert-file")
				tlsKeyFile := c.String("tls-key-file")
				cer, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
				if err != nil {
					return err
				}

				serverConfig.TLSServerConfig = &tls.Config{
					Certificates: []tls.Certificate{cer},
				}
			}

			serverConfig.Addr = fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))
			server, err := server.NewServer(serverConfig, logger)
			if err != nil {
				return err
			}

			err = server.Start()
			if err != nil {
				return err
			}

			return nil
		},
	}
}
