package cli

import (
	"log"

	"github.com/codegangsta/cli"
)

func Server() cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Start the Cimple server",
		Action: func(c *cli.Context) {
			log.Printf("server")
		},
	}
}
