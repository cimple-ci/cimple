package cli

import (
	"log"

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
				Name:  "Host",
				Usage: "The host to bind to",
				Value: "127.0.0.1",
			},
		},
		Action: func(c *cli.Context) {
			logging.SetDefaultLogger("Server", os.Stdout)
			logger := logging.CreateLogger("Server", os.Stdout)

			serverConfig := server.DefaultConfig()
			serverConfig.Addr = fmt.Sprintf("%s:%s", c.String("Host"), c.String("port"))
			server, err := server.NewServer(serverConfig, logger)
			if err != nil {
				log.Fatal(err)
			}

			err = server.Start()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
