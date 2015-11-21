package main

import (
	"github.com/codegangsta/cli"
	cimpleCli "github.com/lukesmith/cimple/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Cimple"
	app.Usage = "Cimple build system"
	app.Commands = []cli.Command{
		cimpleCli.Run(),
		cimpleCli.Server(),
	}

	app.Run(os.Args)
}
