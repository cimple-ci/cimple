package main

import (
	"fmt"
	"os"

	cimpleCli "github.com/lukesmith/cimple/cli/commands"
	"github.com/urfave/cli"
)

var (
	Revision  string
	VERSION   string
	BuildDate string
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "version=%s\n\nrevision=%s\ndate=%s\n", c.App.Version, Revision, BuildDate)
	}

	app := cli.NewApp()
	app.Name = "Cimple"
	app.Usage = "Cimple build system"
	app.Version = VERSION
	app.Commands = []cli.Command{
		cimpleCli.Run(),
		cimpleCli.Server(),
		cimpleCli.Agent(),
		cimpleCli.Config(),
		cimpleCli.Agents(),
		cimpleCli.Builds(),
	}

	app.Run(os.Args)
}
