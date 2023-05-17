package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) AccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "account",
		Usage: "Show account information",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return nil
			}
			ac, err := c.surgesh.Account()
			if err != nil {
				return err
			}
			fmt.Printf("%-6s: %s\n", "Email", ac.Email)
			fmt.Printf("%-6s: %s\n", "ID", ac.ID)
			fmt.Printf("%-6s: %s\n", "UUID", ac.UUID)
			fmt.Printf("%-6s: %s\n", "Plan", ac.Plan.Name)
			fmt.Printf("\n[FEATURES]\n%s\n", strings.Join(ac.Plan.Perks, "\n"))

			return nil
		},
	}
}
