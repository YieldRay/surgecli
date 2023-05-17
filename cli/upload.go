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
)

func (c *privateSurgeCLI) UploadCommand() *cli.Command {
	var isCNAME int

	return &cli.Command{
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
			}},
		Action: func(cCtx *cli.Context) error {
			if email := c.surgesh.Whoami(); email == "" {
				fmt.Println("<YOU ARE NOT LOGGED IN>")
				return nil
			}

			dir := cCtx.Args().Get(0)

			if dir == "" {
				fmt.Println("Usage: surgecli upload <path_to_dir> <domain>")
				fmt.Println()
				fmt.Println(`Use "<CUSTOM_SUBDOMAIN>.surge.sh" if you do not have your own domain`)
				fmt.Println("To setup custom domain, see")
				fmt.Print("https://surge.world/ \nhttps://surge.sh/help/adding-a-custom-domain")
				return nil
			} else {
				if !isDir(dir) {
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
				fmt.Println("Usage: surgecli upload <path_to_dir> <domain>")
				fmt.Println()
				absPath, _ := filepath.Abs(dir)
				fmt.Println("You are going to upload local directory: ", absPath)
				fmt.Printf("You haven't specify a domain and the %s file does not provide a domain\n", cnameFilePath)
				return nil
			}

			if err := c.surgesh.Upload(domain, dir, onUploadEvent); err != nil {
				return err
			} else {
				if isCNAME > 0 {
					return os.WriteFile(cnameFilePath, []byte(domain), 0666)
				}
				return nil
			}
		},
	}
}

func isDir(f string) bool {
	fi, err := os.Stat(f)
	return !os.IsNotExist(err) && fi.IsDir()
}

func onUploadEvent(byteLine []byte) {
	if len(byteLine) == 0 {
		return
	}

	m := make(map[string]any)
	json.Unmarshal(byteLine, &m)

	switch m["type"].(string) {
	case "progress":
		{
			p := types.OnUploadProgress{}
			json.Unmarshal(byteLine, &p)

			if p.End {
				fmt.Printf("%-7s <end>\n", p.Id)
			} else {
				fmt.Printf("%-7s %s \n", p.Id, p.File)
			}
			percentage := fmt.Sprintf("%.2f%%", float32(p.Written)*100/float32(p.Total))
			fmt.Printf("%-7s %d/%d\n", percentage, p.Written, p.Total)
		}
	case "subscription":
		{
		}
	case "ip":
		{
			fmt.Println()
			fmt.Printf("IP %s\n", m["data"].(map[string]any)["ip"])
			fmt.Println()
		}
	case "info":
		{
			p := types.OnUploadInfo{}
			json.Unmarshal(byteLine, &p)

			for _, url := range p.Urls {
				fmt.Printf("[%-12s]\nhttps://%s\n", url.Name, url.Domain)
			}

		}
	}
}
