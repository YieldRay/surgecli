package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) LogoutCommand() *cli.Command {
	return &cli.Command{
		Name:  "logout",
		Usage: "Logout from surge.sh",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return fmt.Errorf("unauthorized")
			}

			email, err := c.surgesh.Logout()
			if err != nil {
				return err
			}
			fmt.Printf("Logout Success as %s\n", email)
			return nil
		},
	}
}
