package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/types"
	"github.com/yieldray/surgecli/utils"
)

func init() {
	var isJSON int

	Commands = append(Commands,
		&cli.Command{
			Name:    "account",
			Aliases: []string{"plan"},
			Usage:   "Show account information",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "json",
					Usage: "Only print JSON",
					Count: &isJSON,
				}},
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

				// only print json
				if isJSON > 0 {
					fmt.Print(utils.JSONStringify(acc))
					return nil
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
