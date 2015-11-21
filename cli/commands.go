package cli

import (
	"github.com/codegangsta/cli"
	"log"
)

func Run() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run Cimple against the current directory",
		Action: func(c *cli.Context) {
			log.Printf("moo")
		},
	}
}

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
