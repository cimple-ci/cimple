package cli

import (
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/runner"
)

func Run() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run Cimple against the current directory",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "task",
				Usage: "a specific task to run. Note that if the task is set to `skip` it will be run.",
			},
		},
		Action: func(c *cli.Context) {
			runner.Run(c.StringSlice("task"))
		},
	}
}
