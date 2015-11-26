package cli

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/lukesmith/cimple/project"
)

func Config() cli.Command {
	return cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Load the config",
		Action: func(c *cli.Context) {
			cfg, err := project.LoadConfig("cimple.hcl")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%+v", cfg)
			log.Printf("config")
		},
	}
}
