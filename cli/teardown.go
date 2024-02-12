package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/types"
)

func init() {
	Commands = append(Commands,
		&cli.Command{
			Name:      "teardown",
			Aliases:   []string{"delete", "rm"},
			Usage:     "Delete site from current account",
			ArgsUsage: "[<domain>]",
			Action: func(cCtx *cli.Context) error {
				if email := surgesh.Whoami(); email == "" {
					fmt.Print("<YOU ARE NOT LOGGED IN>")
					return fmt.Errorf("unauthorized")
				}

				domain := cCtx.Args().First()

				if domain == "" {
					var list types.List
					var err error
					spinner.New().Title("Fetching...").Action(func() {
						list, err = surgesh.List()
					}).Run()

					if err != nil {
						return err
					} else {
						options := make([]huh.Option[string], 0)
						for _, i := range list {
							options = append(options, huh.NewOption(i.Domain, i.Domain))
						}

						err := huh.NewSelect[string]().
							Title("Select one site to delete").
							Value(&domain).Options(options...).Run()

						if err != nil {
							return err
						}

						confirm := false
						huh.NewConfirm().
							Title(fmt.Sprintf("Are you sure to remove %s?", domain)).
							Affirmative("Remove").Negative("Cancel").
							Value(&confirm).Run()

						if !confirm {
							return fmt.Errorf("user aborted")
						}
					}
				}

				var err error
				var teardown types.Teardown
				spinner.New().Title("Removing...").Action(func() {
					teardown, err = surgesh.Teardown(domain)
				}).Run()

				if err != nil {
					return err
				} else {
					fmt.Println(teardown.Msg)
					return nil
				}
			},
		})
}
