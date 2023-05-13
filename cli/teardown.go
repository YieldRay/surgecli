package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) TeardownCommand() *cli.Command {
	return &cli.Command{
		Name:  "teardown",
		Usage: "delete site from surge.sh",
		Action: func(cCtx *cli.Context) error {
			domain := cCtx.Args().Get(0)

			if domain == "" {
				fmt.Println("please specify a domain to teardown")
				return nil
			}

			teardown, err := c.surgesh.Teardown(domain)
			if err != nil {
				return err
			}

			fmt.Println(teardown)

			return nil
		},
	}
}
