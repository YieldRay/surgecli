package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/api"
	"github.com/yieldray/surgecli/types"
	"github.com/yieldray/surgecli/utils"
)

func init() {
	var isCNAME int
	var isSilent int
	var isJSON int
	var version string

	Commands = append(Commands,
		&cli.Command{
			Name:      "upload",
			Aliases:   []string{"deploy"},
			Usage:     "Upload a directory (i.e. deploy a project) to surge.sh",
			ArgsUsage: "<path_to_dir> <domain>",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "CNAME",
					Aliases: []string{"cname"},
					Usage:   "Write the domain to CNAME file",
					Count:   &isCNAME,
				}, &cli.BoolFlag{
					Name:  "silent",
					Usage: "Do not print the uploading info, only print the result",
					Count: &isSilent,
				}, &cli.BoolFlag{
					Name:  "json",
					Usage: "Only print the final info as json",
					Count: &isJSON,
				}, &cli.StringFlag{
					Name:        "version",
					Usage:       "Customize the version string",
					Value:       version,
					DefaultText: api.Version,
				}},
			Action: func(cCtx *cli.Context) error {
				if email := surgesh.Whoami(); email == "" {
					fmt.Println("<YOU ARE NOT LOGGED IN>")
					return fmt.Errorf("unauthorized")
				}

				dir := cCtx.Args().Get(0)

				if dir == "" {
					fmt.Printf("Usage: %s upload <path_to_dir> <domain>\n\n", os.Args[0])
					fmt.Println(`Use "<CUSTOM-SUBDOMAIN>.surge.sh" if you do not have your own domain`)
					fmt.Println("To setup custom domain, see")
					fmt.Print("https://surge.world/ \nhttps://surge.sh/help/adding-a-custom-domain")
					return nil
				}

				domain := cCtx.Args().Get(1)
				cnameFilePath := path.Join(dir, "CNAME")

				if domain == "" {
					domain, _ = utils.ReadTextFileTrim(cnameFilePath)
				}

				if domain == "" {
					fmt.Printf("Usage: %s upload <path_to_dir> <domain>\n\n", os.Args[0])
					absPath, _ := filepath.Abs(dir)
					fmt.Println("You are going to upload local directory: ", absPath)
					fmt.Printf("You haven't specify a domain and the %s file does not provide a domain\n", cnameFilePath)
					return fmt.Errorf("command failed")
				}

				if isSilent == 0 {
					fmt.Fprint(os.Stderr, "Preparing for upload...")
				}

				if err := clearOnError(surgesh.Upload(domain, dir, func(byteLine []byte) {
					onUploadEvent(byteLine, isSilent > 0, isJSON > 0)
				})); err != nil {
					return err
				} else {
					if isCNAME > 0 {
						return clearOnError(os.WriteFile(cnameFilePath, []byte(domain), 0666))
					}
					return nil
				}
			},
		})
}

func clearOnError(err error) error {
	if err != nil {
		utils.ClearLine()
	}
	return err
}

func onUploadEvent(byteLine []byte, isSilent bool, isJSON bool) {
	if len(byteLine) == 0 {
		return
	}

	m := make(map[string]any)
	json.Unmarshal(byteLine, &m)

	pad14 := func(s string) string {
		pad := (14 - len(s)) / 2
		if pad <= 0 {
			return s
		}
		return strings.Repeat(" ", pad) + s
	}

	// print the result info
	handleInfo := func() {
		p := types.OnUploadInfo{}
		json.Unmarshal(byteLine, &p)

		if isJSON {
			fmt.Print(utils.JSONStringify(p))
		} else {
			for _, url := range p.Urls {
				fmt.Printf("[%-14s]\nhttps://%s\n", pad14(url.Name), url.Domain)
			}
		}
	}

	// silent mode, fast return
	if isSilent {
		if m["type"].(string) == "info" {
			handleInfo()
		}
		return // skip print progress
	}

	eprintf := func(format string, a ...any) {
		utils.ClearLineStderr()
		fmt.Fprintf(os.Stderr, format, a...)
	}

	// not silent, print upload progress info
	switch m["type"].(string) {
	case "progress":
		{
			p := types.OnUploadProgress{}
			json.Unmarshal(byteLine, &p)

			percentage := fmt.Sprintf("%.2f%%", float32(p.Written)*100/float32(p.Total))
			if p.End {
				eprintf("%-7s [%-7s %s/%s]", p.Id,
					percentage, utils.FormatBytes(p.Written), utils.FormatBytes(p.Total))
			} else {
				eprintf("%-7s [%-7s %s/%s] %s", p.Id,
					percentage, utils.FormatBytes(p.Written), utils.FormatBytes(p.Total), p.File)
			}

		}
	case "subscription":
		{
		}
	case "ip":
		{
			utils.ClearLine()
			if !isJSON {
				fmt.Printf("[%-14s]\n%s\n", pad14("IP Address"), m["data"].(map[string]any)["ip"])
			}
		}
	case "info":
		{
			handleInfo()
		}
	}
}
