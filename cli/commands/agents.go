package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/api"
	"github.com/olekukonko/tablewriter"
)

func Agents() cli.Command {
	return cli.Command{
		Name:  "agents",
		Usage: "Lists the agents connected to a remote server",
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
			api, err := api.NewApiClient()
			if err != nil {
				log.Fatal(err)
			}

			api.ServerUrl = "http://" + c.String("server-addr") + ":" + c.String("server-port")

			agents, err := api.GetAgents()
			if err != nil {
				log.Fatal(err)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"AGENT ID", "CONNECTED", "BUSY"})

			for _, agent := range agents {
				table.Append([]string{
					fmt.Sprintf("%s", agent.Id),
					fmt.Sprintf("%s", agent.ConnectedDate),
					fmt.Sprintf("%t", agent.Busy)})
			}

			table.Render()
		},
	}
}
