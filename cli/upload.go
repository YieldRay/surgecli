package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/types"
)

func (c *privateSurgeCLI) UploadCommand() *cli.Command {
	return &cli.Command{
		Name:    "upload",
		Aliases: []string{"deploy"},
		Usage:   "upload directory to surge.sh with specified domain",
		Action: func(cCtx *cli.Context) error {

			dir := cCtx.Args().Get(0)

			if dir == "" {
				dir = "./"
			}

			domain := cCtx.Args().Get(1)

			if domain == "" {
				b, _ := os.ReadFile(path.Join(dir, "CNAME"))
				domain = strings.Trim(string(b), " ")
			}
			if domain == "" {
				fmt.Println(`Usage: surgecli upload <path_to_dir> <domain>`)
				fmt.Println(`domain cannot be empty!`)
				fmt.Println(`use "<CUSTOM_SUB_DOMAIN>.surge.sh" if you do not have a domain`)
				fmt.Println("to setup custom domain, see https://surge.sh/help/adding-a-custom-domain")
				return nil
			} else {
				return c.surgesh.Upload(domain, dir, onUploadEvent)
			}

		},
	}
}

func onUploadEvent(byteLine []byte) {
	m := make(map[string]interface{})
	json.Unmarshal(byteLine, &m)

	t := m["type"].(string)
	switch t {
	case "progress":
		{
			p := types.OnUploadProgress{}
			json.Unmarshal(byteLine, &p)

			if p.End {
				fmt.Printf("%-6s <end>\n", p.Id)
			} else {
				fmt.Printf("%-6s %s \n", p.Id, p.File)
			}
			fmt.Printf("%.2f%% %d/%d\n", float32(p.Written)*100/float32(p.Total), p.Written, p.Total)
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
				fmt.Printf("%s\nhttps://%s\n", url.Name, url.Domain)
			}

		}
	}
}
