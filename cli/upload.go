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
	return &cli.Command{
		Name:      "upload",
		Aliases:   []string{"deploy"},
		Usage:     "Upload a directory (a.k.a. deploy a project) to surge.sh",
		ArgsUsage: "<path_to_dir> <domain>",
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
				fmt.Println("To setup custom domain, see ")
				fmt.Println("https://surge.world/")
				fmt.Println("https://surge.sh/help/adding-a-custom-domain")
				return nil
			} else {
				if !isDir(dir) {
					return fmt.Errorf("%s is not a directory", dir)
				}
			}

			domain := cCtx.Args().Get(1)

			if domain == "" {
				b, _ := os.ReadFile(path.Join(dir, "CNAME"))
				domain = strings.Trim(string(b), " ")
			}

			if domain == "" {
				fmt.Println("Usage: surgecli upload <path_to_dir> <domain>")
				fmt.Println()
				absPath, _ := filepath.Abs(dir)
				fmt.Println("You are going to upload local directory: ", absPath)
				fmt.Println("You have NOT specify a domain, please enter a domain")
				fmt.Println(`Use "<CUSTOM_SUBDOMAIN>.surge.sh" if you do not have your own domain`)
				fmt.Println("To setup custom domain, see https://surge.sh/help/adding-a-custom-domain")
				fmt.Println()
				fmt.Print("Please Enter Your Domain (Ctrl+C To Quit): ")
				fmt.Scanf("%s\n", &domain)
				if len(domain) < 3 {
					return fmt.Errorf("domain is invalid")
				}

				cnameFilePath := path.Join(dir, "CNAME")
				fmt.Println()
				fmt.Printf("Do you want to write the domain to %s file?\n", cnameFilePath)
				fmt.Print("This will allow the cli to remember that domain (yes/no): ")
				confirmText := ""
				fmt.Scanf("%s\n", &confirmText)
				if confirmText == "yes" {
					os.WriteFile(cnameFilePath, []byte(domain), 0666)
				}
			}

			fmt.Println()
			return c.surgesh.Upload(domain, dir, onUploadEvent)
		},
	}
}

func isDir(f string) bool {
	if fi, err := os.Stat(f); err != nil {
		return false
	} else {
		return fi.IsDir()
	}
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
