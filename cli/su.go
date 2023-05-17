package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	surgeUtils "github.com/yieldray/surgecli/utils"
)

func (c *privateSurgeCLI) SuCommand() *cli.Command {
	return &cli.Command{
		Name:      "su",
		Usage:     "switch user",
		ArgsUsage: "<email>",
		Action: func(cCtx *cli.Context) error {
			email := cCtx.Args().First()
			if len(email) == 0 {
				fmt.Printf("Usage: %s su <email>\n\n", os.Args[0])
				emails := surgeUtils.ConfGetEmailList()
				if len(emails) == 0 {
					fmt.Println("there is no user stored in config file")
				} else {
					fmt.Printf("email list:\n")
					fmt.Println(strings.Join(emails, "\n"))
				}
			}
			if err := surgeUtils.ConfUseEmail(email); err != nil {
				return err
			} else {
				fmt.Printf("Switched to user %s\n", email)
				return nil
			}
		},
	}
}
