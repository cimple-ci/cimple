package main

import (
	"os"

	"github.com/codegangsta/cli"
	cimpleCli "github.com/lukesmith/cimple/cli/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "Cimple"
	app.Usage = "Cimple build system"
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
