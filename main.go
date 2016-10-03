package main

import (
	"os"

	cimpleCli "github.com/lukesmith/cimple/cli/commands"
	"github.com/urfave/cli"
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
