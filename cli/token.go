package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/utils"
)

func init() {

	Commands = append(Commands,
		&cli.Command{
			Name:      "token",
			Usage:     "Show current token",
			ArgsUsage: "[<username> <password>]",
			Action: func(cCtx *cli.Context) error {
				if _, token, err := utils.ReadNetrc(); err != nil {
					return err
				} else {
					fmt.Println(token)
					return nil
				}
			},
		})
}
