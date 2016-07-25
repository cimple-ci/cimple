package cli

import (
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/api"
	"log"
)

func Builds() cli.Command {
	return cli.Command{
		Name:  "builds",
		Usage: "Provides access to managing builds",
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
					client, err := api.NewApiClient()
					if err != nil {
						log.Fatal(err)
					}

					client.ServerUrl = "http://" + c.String("server-addr") + ":" + c.String("server-port")

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
		},
	}
}
