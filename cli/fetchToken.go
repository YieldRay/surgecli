package cli

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/api"
)

func (c *privateSurgeCLI) FetchTokenCommand() *cli.Command {
	var isLocal int

	return &cli.Command{
		Name:      "fetch-token",
		Usage:     "Fetch token by email and password, but do not save the token like login command",
		ArgsUsage: "<username> <password>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "local",
				Usage: "Only print token from ~/.netrc file, rather than login to server",
				Count: &isLocal,
			}},
		Action: func(cCtx *cli.Context) error {

			// print local token
			if isLocal > 0 {
				if _, token, err := api.ReadNetrc(); err != nil {
					return err
				} else {
					fmt.Print(token)
					return nil
				}
			}

			// otherwise, login to server and print token
			username := cCtx.Args().Get(0)
			password := cCtx.Args().Get(1)

			if password == "" {
				fmt.Print("Usage: surgecli fetch-token <username> <password>")
				return nil
			}

			if tokens, err := api.Token(http.DefaultClient, username, password); err != nil {
				return err
			} else {
				fmt.Print(tokens.Token)
				return nil
			}

		},
	}
}
