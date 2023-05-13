package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) TeardownCommand() *cli.Command {
	return &cli.Command{
		Name:      "teardown",
		Usage:     "Delete site from surge.sh",
		ArgsUsage: "<domain>",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return nil
			}

			domain := cCtx.Args().First()

			if domain == "" {
				fmt.Println("Usage: surgecli teardown <domain>")
				fmt.Println("please specify a domain to teardown")
				return nil
			}

			teardown, err := c.surgesh.Teardown(domain)
			if err != nil {
				return err
			}

			fmt.Println(teardown.Msg)
			return nil
		},
	}
}
