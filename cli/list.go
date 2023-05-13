package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List my sites",
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return nil
			}

			if list, err := c.surgesh.List(); err != nil {
				return err
			} else {
				for _, site := range list {
					fmt.Printf("%-40s [%s]\n", site.Domain, site.TimeAgoInWords)
					fmt.Printf("%-7s: https://%s\n", "URL", site.Domain)
					fmt.Printf("%-7s: https://%s\n", "Preview", site.Preview)
					fmt.Printf("v%-7s %6d file(s)  %6d byte(s) \n", site.CliVersion, site.PublicFileCount, site.PublicTotalSize)
					// fmt.Printf("%s - %s\n", time.Unix(site.UploadStartTime/1000, 0), time.Unix(site.UploadEndTime/1000, 0))
					fmt.Println()
				}
			}
			return nil
		},
	}
}
