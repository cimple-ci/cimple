package cli

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/server"
	"os"
)

func Server() cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Start the Cimple server",
		Action: func(c *cli.Context) {
			log.Printf("server")

			serverConfig := server.DefaultConfig()
			serverConfig.Addr = ":8080"
			server, err := server.NewServer(serverConfig, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}

			go func() {
				err = server.Start()
				if err != nil {
					log.Fatal(err)
				}
			}()
		},
	}
}
