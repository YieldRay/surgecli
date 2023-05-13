package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) LogoutCommand() *cli.Command {
	return &cli.Command{
		Name:  "logout",
		Usage: "logout to surge.sh",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return nil
			}
			email, err := c.surgesh.Logout()
			if err == nil {
				fmt.Printf("Logout Success as %s\n", email)
			}
			return err
		},
	}
}
