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
		Name:  "surgecli",
		Usage: "thrid party surge.sh cli",
		Commands: []*cli.Command{
			surgeCLI.LoginCommand(),
			surgeCLI.LogoutCommand(),
			surgeCLI.WhoamiCommand(),
			surgeCLI.AccountCommand(),
			surgeCLI.ListCommand(),
			surgeCLI.TeardownCommand(),
			surgeCLI.UploadCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api",
				Value:       "https://surge.surge.sh",
				Usage:       "configure the API host",
				Destination: &surgeCLI.API_HOST,
			},
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
