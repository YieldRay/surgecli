package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	surgecli "github.com/yieldray/surgecli/cli"
)

func main() {
	surgeCLI := surgecli.SurgeCLI

	app := &cli.App{
		Name:     "surgecli",
		Usage:    "thrid party surge.sh cli",
		Commands: surgeCLI.Commands(),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "toggle debug on",
				Count: &surgeCLI.DEBUG,
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
