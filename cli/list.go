package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list my sites",
		Action: func(cCtx *cli.Context) error {
			list, err := c.surgesh.List()
			if err != nil {
				return err
			}
			for _, site := range list {
				fmt.Println(site)
			}
			return nil
		},
	}
}
