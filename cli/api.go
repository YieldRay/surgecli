package cli

import (
	"fmt"
	"net/url"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/utils"
)

func init() {
	var isShow int
	var isReset int

	Commands = append(Commands,
		&cli.Command{
			Name:      "api",
			Usage:     "Set or Show API base URL, the official one is https://surge.surge.sh",
			ArgsUsage: "[<baseURL>]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "show",
					Usage: "Show current API base URL",
					Count: &isShow,
				}, &cli.BoolFlag{
					Name:  "reset",
					Usage: "Reset API base URL to official one",
					Count: &isReset,
				}},
			Action: func(cCtx *cli.Context) error {

				prevApi := utils.ConfGetApi()
				if isShow > 0 {
					fmt.Println(prevApi)
					return nil
				}

				api := cCtx.Args().First()

				if isReset > 0 {
					api = "https://surge.surge.sh"
				} else if len(api) == 0 {
					huh.NewInput().Title("Enter new API base").Placeholder(prevApi).
						Suggestions([]string{"https://surge.surge.sh"}).
						Validate(func(s string) error {
							_, err := url.ParseRequestURI(s)
							return err
						}).Value(&api).Run()
				}

				if err := utils.ConfSetApi(api); err != nil {
					return err
				} else {
					fmt.Printf("API base URL has been set to %s\n", api)
					return nil
				}
			},
		})
}
