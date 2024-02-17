package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/api"
	"github.com/yieldray/surgecli/types"
)

func init() {

	Commands = append(Commands,
		&cli.Command{
			Name:      "fetch-token",
			Usage:     "Fetch token by email and password, but do not save the token like login command",
			ArgsUsage: "[<username> <password>]",
			Action: func(cCtx *cli.Context) error {

				// otherwise, login to server and print token
				username := cCtx.Args().Get(0)
				password := cCtx.Args().Get(1)

				if password == "" {
					huh.NewInput().Title("Enter you username (email)").Value(&username).Run()
					huh.NewInput().Title("Enter you password").Value(&password).Run()
				}

				var tokens types.Token
				var err error
				spinner.New().Title("Fetching...").Action(func() {
					tokens, err = api.Token(httpClient, username, password)
				}).Run()

				if err != nil {
					return err
				} else {
					fmt.Println(tokens.Token)
					return nil
				}
			},
		})
}
