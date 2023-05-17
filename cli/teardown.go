package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) TeardownCommand() *cli.Command {
	return &cli.Command{
		Name:      "teardown",
		Aliases:   []string{"delete"},
		Usage:     "Delete site from surge.sh",
		ArgsUsage: "<domain>",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Print("<YOU ARE NOT LOGGED IN>")
				return nil
			}

			domain := cCtx.Args().First()

			if domain == "" {
				fmt.Printf("Usage: %s teardown <domain>\n", os.Args[0])
				fmt.Println("Please specify a domain to teardown")
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
