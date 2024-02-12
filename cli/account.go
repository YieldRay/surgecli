package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/types"
)

func init() {
	Commands = append(Commands,
		&cli.Command{
			Name:    "account",
			Aliases: []string{"me"},
			Usage:   "Show account information",
			Action: func(cCtx *cli.Context) error {

				if email := surgesh.Whoami(); email == "" {
					fmt.Println("<YOU ARE NOT LOGGED IN>")
					return fmt.Errorf("unauthorized")
				}

				var acc types.Account
				var err error
				spinner.New().Title("Fetching...").Action(func() {
					acc, err = surgesh.Account()
				}).Run()

				if err != nil {
					return err
				}

				fmt.Printf("%-6s: %s\n", "Email", acc.Email)
				fmt.Printf("%-6s: %s\n", "ID", acc.ID)
				fmt.Printf("%-6s: %s\n", "UUID", acc.UUID)
				fmt.Printf("%-6s: %s\n", "Plan", acc.Plan.Name)
				fmt.Printf("\n[FEATURES]\n%s\n", strings.Join(acc.Plan.Perks, "\n"))
				return nil
			},
		})
}
