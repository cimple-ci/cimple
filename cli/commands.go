package cli

import (
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/project"
	"log"
	"os"
)

func Run() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run Cimple against the current directory",
		Action: func(c *cli.Context) {
			cfg, err := project.LoadConfig("cimple.hcl")
			if err != nil {
				log.Fatal(err)
			}

			project.Run(cfg, os.Stdout, os.Stderr)
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
