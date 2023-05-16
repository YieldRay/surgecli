package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func (c *privateSurgeCLI) ListCommand() *cli.Command {
	var isShort int

	return &cli.Command{
		Name:  "list",
		Usage: "List my sites",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "short",
				Usage: "Only print domain list",
				Count: &isShort,
			}},
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Print("<YOU ARE NOT LOGGED IN>")
				return nil
			}

			if list, err := c.surgesh.List(); err != nil {
				return err
			} else {

				// only print domain
				if isShort > 0 {
					domains := make([]string, 0, len(list))
					for _, site := range list {
						domains = append(domains, site.Domain)
					}
					fmt.Print(strings.Join(domains, "\n"))
					return nil
				}

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
