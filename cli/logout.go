package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands,
		&cli.Command{
			Name:  "logout",
			Usage: "Logout from surge.sh",
			Action: func(cCtx *cli.Context) error {
				if email := surgesh.Whoami(); email == "" {
					fmt.Println("<YOU ARE NOT LOGGED IN>")
					return fmt.Errorf("unauthorized")
				}

				email, err := surgesh.Logout()
				if err != nil {
					return err
				} else {
					fmt.Printf("Logout Success as %s\n", email)
					return nil
				}
			},
		})
}
