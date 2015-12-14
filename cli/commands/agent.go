package cli

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/agent"
	"os"
)

func Agent() cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Start the Cimple agent",
		Action: func(c *cli.Context) {
			log.Printf("agent")

			agentConfig, err := agent.DefaultConfig()
			agentConfig.ServerPort = "8080"
			agent, err := agent.NewAgent(agentConfig, os.Stdout)

			go func() {
				err = agent.Start()
				if err != nil {
					log.Fatal(err)
				}
			}()
		},
	}
}
