package cli

import (
	"fmt"
	"os"

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
				fmt.Printf("Usage: %s login <username> <password>\n", os.Args[0])
				return fmt.Errorf("command failed")
			}

			email, err := c.surgesh.Login(username, password)
			if err == nil {
				fmt.Printf("Login Success as %s\n", email)
			}
			return err
		},
	}
}
