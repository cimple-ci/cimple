package cli

import (
	"github.com/codegangsta/cli"
	"log"
	//"github.com/lukesmith/cimple/api"
)

func Builds() cli.Command {
	return cli.Command{
		Name:  "builds",
		Usage: "Provides access to managing builds",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "server-addr",
				Usage: "The Cimple server address.",
				Value: "localhost",
			},
			cli.StringFlag{
				Name:  "server-port",
				Usage: "The Cimple server port.",
				Value: "8080",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:  "submit",
				Usage: "Submits a new build to the server",
				Action: func(c *cli.Context) {
					log.Print("Submit")
				},
			},
		},
		Action: func(c *cli.Context) {
		},
	}
}
