package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) LoginCommand() *cli.Command {
	return &cli.Command{
		Name:      "login",
		Usage:     "Login (or create new account) to surge.sh",
		ArgsUsage: "<username> <password>",
		Action: func(cCtx *cli.Context) error {
			username := cCtx.Args().Get(0)
			password := cCtx.Args().Get(1)

			if password == "" {
				fmt.Print("Usage: surgecli login <username> <password>")
				return nil
			}

			email, err := c.surgesh.Login(username, password)
			if err == nil {
				fmt.Printf("Login Success as %s", email)
			}
			return err
		},
	}
}
