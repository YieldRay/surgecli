package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/types"
	"github.com/yieldray/surgecli/utils"
)

func init() {
	var isCNAME int
	var isSilent int

	Commands = append(Commands,
		&cli.Command{
			Name:      "upload",
			Aliases:   []string{"deploy"},
			Usage:     "Upload a directory (a.k.a. deploy a project) to surge.sh",
			ArgsUsage: "<path_to_dir> <domain>",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "CNAME",
					Aliases: []string{"cname"},
					Usage:   "write the domain to CNAME file",
					Count:   &isCNAME,
				}, &cli.BoolFlag{
					Name:  "silent",
					Usage: "do not print the uploading info, only print the result",
					Count: &isSilent,
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
				} else {
					if !utils.IsDir(dir) {
						return fmt.Errorf("%s is not a directory", dir)
					}
				}

				domain := cCtx.Args().Get(1)
				cnameFilePath := path.Join(dir, "CNAME")

				if domain == "" {
					b, _ := os.ReadFile(cnameFilePath)
					domain = strings.Trim(string(b), " ")
				}

				if domain == "" {
					fmt.Printf("Usage: %s upload <path_to_dir> <domain>\n\n", os.Args[0])
					absPath, _ := filepath.Abs(dir)
					fmt.Println("You are going to upload local directory: ", absPath)
					fmt.Printf("You haven't specify a domain and the %s file does not provide a domain\n", cnameFilePath)
					return fmt.Errorf("command failed")
				}

				if isSilent == 0 {
					fmt.Print("Preparing for upload...")
				}

				if err := surgesh.Upload(domain, dir, func(byteLine []byte) {
					onUploadEvent(byteLine, isSilent > 0)
				}); err != nil {
					return err
				} else {
					if isCNAME > 0 {
						return os.WriteFile(cnameFilePath, []byte(domain), 0666)
					}
					return nil
				}
			},
		})
}

func onUploadEvent(byteLine []byte, isSilent bool) {
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

		for _, url := range p.Urls {
			fmt.Printf("[%-14s]\nhttps://%s\n", pad14(url.Name), url.Domain)
		}
	}

	// silent mode, fast return
	if isSilent {
		if m["type"].(string) == "info" {
			handleInfo()
		}
		return
	}

	printf := func(format string, a ...any) {
		utils.ClearLine()
		fmt.Printf(format, a...)
	}

	// not silent, print upload progress info
	switch m["type"].(string) {
	case "progress":
		{
			p := types.OnUploadProgress{}
			json.Unmarshal(byteLine, &p)

			percentage := fmt.Sprintf("%.2f%%", float32(p.Written)*100/float32(p.Total))
			if p.End {
				printf("%-7s [%-7s %s/%s]", p.Id,
					percentage, utils.FormatBytes(p.Written), utils.FormatBytes(p.Total))
			} else {
				printf("%-7s [%-7s %s/%s] %s", p.Id,
					percentage, utils.FormatBytes(p.Written), utils.FormatBytes(p.Total), p.File)
			}

		}
	case "subscription":
		{
		}
	case "ip":
		{
			utils.ClearLine()
			fmt.Printf("[%-14s]\n%s\n", pad14("IP Address"), m["data"].(map[string]any)["ip"])
		}
	case "info":
		{
			handleInfo()
		}
	}
}
