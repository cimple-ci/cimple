package cli

import (
	"fmt"
	"github.com/cimple-ci/cimple-go-api"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"log"
	"os"
)

func newApiClient(c *cli.Context) (*api.ApiClient, error) {
	client, err := api.NewApiClient()
	if err != nil {
		return nil, err
	}

	client.ServerUrl = "http://" + c.String("server-addr") + ":" + c.String("server-port")

	return client, nil
}

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
					client, err := newApiClient(c)
					if err != nil {
						log.Fatal(err)
					}

					buildOptions := &api.BuildSubmissionOptions{
						Url:    c.Args().Get(0),
						Commit: c.Args().Get(1),
					}
					err = client.SubmitBuild(buildOptions)
					if err != nil {
						log.Fatal(err)
					}
				},
			},
		},
		Action: func(c *cli.Context) {
			client, err := newApiClient(c)
			if err != nil {
				log.Fatal(err)
			}

			builds, err := client.ListBuilds()
			if err != nil {
				log.Fatal(err)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "SUBMISSION DATE"})

			for _, build := range builds {
				table.Append([]string{
					fmt.Sprintf("%s", build.Id),
					fmt.Sprintf("%s", build.SubmissionDate),
				})
			}

			table.Render()
		},
	}
}
