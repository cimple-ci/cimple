package cli

import (
	"fmt"
	"github.com/lukesmith/cimple/runner"
	"github.com/urfave/cli"
	"strings"
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
			cli.StringSliceFlag{
				Name:  "secret",
				Usage: "specifies a `SECRET` to make available to the tasks. Secrets must be defined in the format `type:key:password`",
			},
		},
		Action: func(c *cli.Context) error {
			ss, err := makeCliSecretStore(c.StringSlice("secret"))
			if err != nil {
				return err
			}

			runOptions := &runner.RunOptions{
				Journal: &runner.JournalSettings{
					Driver: c.String("journal-driver"),
					Format: c.String("journal-format"),
				},
				Context: c.String("run-context"),
				Secrets: ss,
			}

			return runner.Run(runOptions, c.StringSlice("task"))
		},
	}
}

func makeCliSecretStore(vals []string) (*cliSecretStore, error) {
	ss := &cliSecretStore{}
	ss.secrets = make(map[string]map[string]string)

	for _, s := range vals {
		parts := strings.Split(s, ":")

		if len(parts) != 3 {
			return nil, fmt.Errorf("Secret defined without 3 parts. Ensure the secret is in the format `type:key:password`")
		}

		if _, ok := ss.secrets[parts[0]]; !ok {
			ss.secrets[parts[0]] = make(map[string]string)
		}

		ss.secrets[parts[0]][parts[1]] = parts[2]
	}

	return ss, nil
}

type cliSecretStore struct {
	secrets map[string]map[string]string
}

func (sv cliSecretStore) Get(secretType string, key string) (string, error) {
	if st, ok := sv.secrets[secretType]; ok {
		if val, ok := st[key]; ok {
			return val, nil
		}
	}

	return "", fmt.Errorf("Failed to find %s password for key %s", secretType, key)
}
