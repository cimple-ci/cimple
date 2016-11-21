package cli

import (
	"github.com/lukesmith/cimple/runner"
	"github.com/urfave/cli"
)

func Run() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run Cimple against the current directory",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "task",
				Usage: "a specific `TASK` to run. Note that if the task is set to `skip` it will be run.",
			},
			cli.StringFlag{
				Name:  "journal-driver",
				Usage: "specifiy the `DRIVER` to send journal messages to. Available options \"console\"",
				Value: "console",
			},
			cli.StringFlag{
				Name:  "journal-format",
				Usage: "specify the output `FORMAT`. Available options \"text\", \"json\"",
				Value: "text",
			},
			cli.StringFlag{
				Name:  "run-context",
				Usage: "specify the `CONTEXT` in which to execute the run in. Available options \"local\", \"server\"",
				Value: "local",
			},
		},
		Action: func(c *cli.Context) {
			runOptions := &runner.RunOptions{
				Journal: &runner.JournalSettings{
					Driver: c.String("journal-driver"),
					Format: c.String("journal-format"),
				},
				Context: c.String("run-context"),
			}
			runner.Run(runOptions, c.StringSlice("task"))
		},
	}
}
