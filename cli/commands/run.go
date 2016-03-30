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
			cli.StringFlag{
				Name:  "syslog",
				Usage: "a Syslog host to send logs to.",
			},
		},
		Action: func(c *cli.Context) {
			runOptions := &runner.RunOptions{
				LogServer: c.String("syslog"),
			}
			runner.Run(runOptions, c.StringSlice("task"))
		},
	}
}
