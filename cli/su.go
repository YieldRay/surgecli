package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	surgeUtils "github.com/yieldray/surgecli/utils"
)

func selectEmail(emails []string, input string) string {
	for _, email := range emails {
		if strings.HasPrefix(email, input) {
			return email
		}
	}
	return input
}

func (c *privateSurgeCLI) SuCommand() *cli.Command {
	return &cli.Command{
		Name:      "su",
		Usage:     "Switch user",
		ArgsUsage: "<email>",
		Action: func(cCtx *cli.Context) error {
			email := cCtx.Args().First()
			emails := surgeUtils.ConfGetEmailList()
			if len(email) == 0 {
				fmt.Printf("Usage: %s su <email>\n\n", os.Args[0])
				if len(emails) == 0 {
					fmt.Println("there is no user stored in config file")
				} else {
					fmt.Printf("[email list]\n")
					fmt.Println(strings.Join(emails, "\n"))
				}
				return nil
			}

			// the final selected email address string
			selection := selectEmail(emails, email)
			if err := surgeUtils.ConfUseEmail(selection); err != nil {
				return err
			} else {
				fmt.Printf("Switched to user %s\n", selection)
				return nil
			}
		},
	}
}
