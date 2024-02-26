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
	var isJSON int

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
				}, &cli.BoolFlag{
					Name:  "json",
					Usage: "Only print JSON",
					Count: &isJSON,
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

					// only print json
					if isJSON > 0 {
						fmt.Print(utils.JSONStringify(list))
						return nil
					}

					// only print domain
					if isShort > 0 {
						for _, site := range list {
							fmt.Println(site.Domain)
						}
						return nil
					}

					// print full info
					s := func(count int) string {
						if count == 1 {
							return fmt.Sprintf("%d file", count)
						}
						return fmt.Sprintf("%d files", count)
					}

					for _, site := range list {
						fmt.Println(site.Domain)
						fmt.Printf("%-7s: https://%s\n", "URL", site.Domain)
						fmt.Printf("%-7s: https://%s\n", "Preview", site.Preview)
						fmt.Printf("v%-7s %-24s %-15s\n", site.CliVersion,
							fmt.Sprintf("%s - %s", s(site.PublicFileCount), utils.FormatBytes(site.PublicTotalSize)),
							fmt.Sprintf("[%s]", site.TimeAgoInWords))
						// fmt.Printf("%s - %s\n", time.Unix(site.UploadStartTime/1000, 0), time.Unix(site.UploadEndTime/1000, 0))
						fmt.Println()
					}
				}
				return nil
			},
		})
}
