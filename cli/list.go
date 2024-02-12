package cli

import (
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/types"
	"github.com/yieldray/surgecli/utils"
)

func init() {
	var isShort int

	Commands = append(Commands,
		&cli.Command{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List my sites",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "short",
					Usage: "Only print domain list",
					Count: &isShort,
				}},
			Action: func(cCtx *cli.Context) error {
				if email := surgesh.Whoami(); email == "" {
					fmt.Println("<YOU ARE NOT LOGGED IN>")
					return fmt.Errorf("unauthorized")
				}

				var list types.List
				var err error
				spinner.New().Title("Fetching...").Action(func() {
					list, err = surgesh.List()
				}).Run()

				if err != nil {
					return err
				} else {
					// only print domain
					if isShort > 0 {
						for _, site := range list {
							fmt.Println(site.Domain)
						}
						return nil
					}

					// print full info
					s := func(fc int) string {
						if fc == 1 {
							return fmt.Sprintf("%d file", fc)
						}
						return fmt.Sprintf("%d files", fc)
					}

					for _, site := range list {
						fmt.Printf("%-40s [%s]\n", site.Domain, site.TimeAgoInWords)
						fmt.Printf("%-7s: https://%s\n", "URL", site.Domain)
						fmt.Printf("%-7s: https://%s\n", "Preview", site.Preview)
						fmt.Printf("v%-7s %-11s %s\n", site.CliVersion, s(site.PublicFileCount), utils.FormatBytes(site.PublicTotalSize))
						// fmt.Printf("%s - %s\n", time.Unix(site.UploadStartTime/1000, 0), time.Unix(site.UploadEndTime/1000, 0))
						fmt.Println()
					}
				}
				return nil
			},
		})
}
