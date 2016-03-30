package cli

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/agent"
	"os"
)

func Agent() cli.Command {
	return cli.Command{
		Name:    "agent",
		Aliases: []string{"a"},
		Usage:   "Start the Cimple agent",
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
		Action: func(c *cli.Context) {
			log.Printf("agent")

			agentConfig, err := agent.DefaultConfig()
			agentConfig.ServerAddr = c.String("server-addr")
			agentConfig.ServerPort = c.String("server-port")
			agent, err := agent.NewAgent(agentConfig, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}

			err = agent.Start()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
