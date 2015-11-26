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
		cimpleCli.Config(),
	}

	app.Run(os.Args)
}
