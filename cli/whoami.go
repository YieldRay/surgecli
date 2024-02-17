package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands,
		&cli.Command{
			Name:  "whoami",
			Usage: "Show who you are logged in as",
			Action: func(cCtx *cli.Context) error {
				if email := surgesh.Whoami(); email == "" {
					fmt.Println("<YOU ARE NOT LOGGED IN>")
					return fmt.Errorf("unauthorized")
				} else {
					fmt.Println(email)
					return nil
				}
			},
		})
}
