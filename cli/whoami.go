package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) WhoamiCommand() *cli.Command {
	return &cli.Command{
		Name:  "whoami",
		Usage: "Show my email",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
			} else {
				fmt.Println(email)
			}
			return nil
		},
	}
}
