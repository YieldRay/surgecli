package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	surgecli "github.com/yieldray/surgecli/cli"
)

func main() {

	app := &cli.App{
		Name:     "surgecli",
		Usage:    "third party surge.sh cli",
		Commands: surgecli.Commands,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "toggle debug on",
				Count: &surgecli.FLAG_DEBUG,
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
