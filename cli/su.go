package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/utils"
)

func init() {
	Commands = append(Commands,
		&cli.Command{
			Name:      "su",
			Usage:     "Switch user",
			ArgsUsage: "[<email>]",
			Action: func(cCtx *cli.Context) error {
				email := cCtx.Args().First()
				emails := utils.ConfGetEmailList()

				if len(emails) == 0 {
					return fmt.Errorf("there is no user stored in config file")
				}

				if len(email) == 0 {
					options := make([]huh.Option[string], 0)
					for _, e := range emails {
						options = append(options, huh.NewOption(e, e))
					}

					err := huh.NewSelect[string]().
						Title("Switch to another account").
						Value(&email).Options(options...).Run()

					if err != nil {
						return err
					}
				}

				if len(email) == 0 {
					return fmt.Errorf("command failed")
				}

				if err := utils.ConfUseEmail(email); err != nil {
					return err
				} else {
					fmt.Printf("Switched to user %s\n", email)
					return nil
				}
			},
		})
}
