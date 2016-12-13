package cli

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gyuho/goraph"
	"github.com/urfave/cli"

	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/project"
)

func Config() cli.Command {
	return cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Load the config",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "configuration",
				Usage: "specify the CONFIGURATION file to load",
				Value: "cimple.hcl",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "Specify the FORMAT to write the config",
				Value: "graphviz",
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("configuration")
			cfg, err := project.LoadConfig(configFile)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Unable to load configuration file - %s - %s", configFile, err), CONFIGURATION_ERROR_CODE)
			}

			if c.String("format") == "graphviz" {
				mt := []build.TaskNode{}
				for _, t := range cfg.Tasks {
					mt = append(mt, t)
				}

				graph := build.PopulateGraph(mt)

				output, err := writeGraphViz(graph)
				os.Stdout.Write([]byte(output))
				if err != nil {
					fmt.Fprintf(os.Stderr, "%+v", err)
					return cli.NewExitError("Unable to generate graphviz representation", CONFIGURATION_ERROR_CODE)
				}
			}

			return nil
		},
	}
}

func writeGraphViz(g goraph.Graph) (string, error) {
	writer := &bytes.Buffer{}
	writer.Write([]byte("digraph build {"))
	writer.Write([]byte("\n"))
	written := make(map[string]bool)

	a := build.Entrypoints(g)

	for _, ep := range a {
		writer.Write([]byte(fmt.Sprintf("\t%s", ep.ID())))

		err := traverse(g, ep, writer, written)
		if err != nil {
			return "", err
		}

		writer.Write([]byte(";\n"))
	}

	writer.Write([]byte("}"))

	return writer.String(), nil
}

func traverse(g goraph.Graph, node goraph.Node, writer *bytes.Buffer, written map[string]bool) error {
	targets, err := g.GetTargets(node.ID())
	if err != nil {
		return err
	}

	for _, tg := range targets {
		writer.Write([]byte(fmt.Sprintf(" -> %s", tg.ID())))
		if _, ok := written[tg.String()]; !ok {
			written[tg.String()] = true
			err := traverse(g, tg, writer, written)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
