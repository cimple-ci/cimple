package cli

import (
	"log"

	"fmt"
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/server"
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
		},
		Action: func(c *cli.Context) {
			serverConfig := server.DefaultConfig()
			serverConfig.Addr = fmt.Sprintf(":%s", c.String("port"))
			server, err := server.NewServer(serverConfig, os.Stdout)
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
