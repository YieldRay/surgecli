package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) AccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "account",
		Usage: "show account information",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return nil
			}
			ac, err := c.surgesh.Account()
			if err != nil {
				return err
			}
			fmt.Println(ac)

			return nil
		},
	}
}
