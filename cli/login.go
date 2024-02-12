package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands,
		&cli.Command{
			Name:      "login",
			Usage:     "Login (or create new account) to surge.sh",
			ArgsUsage: "[<username> <password>]",
			Action: func(cCtx *cli.Context) error {
				username := cCtx.Args().Get(0)
				password := cCtx.Args().Get(1)

				if password == "" {
					huh.NewInput().Title("Enter you username (email)").Value(&username).Run()
					huh.NewInput().Title("Enter you password").Value(&password).Run()
				}

				var email string
				var err error
				spinner.New().Title("Logging in...").Action(func() {
					email, err = surgesh.Login(username, password)
				}).Run()

				if err != nil {
					return err
				} else {
					fmt.Printf("Login Success as %s\n", email)
					return nil
				}
			},
		})
}
