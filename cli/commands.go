package cli

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/agent"
	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/server"
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
				panic(err)
			}

			build, err := build.NewBuild(cfg)
			if err != nil {
				log.Fatal(err)
			}
			log.Print(build)

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

			agentConfig, err := agent.DefaultConfig()
			agentConfig.ServerPort = "8080"
			agent, err := agent.NewAgent(agentConfig, os.Stdout)

			err = agent.Start()
			if err != nil {
				log.Fatal(err)
			}
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
