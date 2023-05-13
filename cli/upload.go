package cli

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

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
				fmt.Println(`Usage: surgecli upload <path_to_dir> <domain>`)
				fmt.Println(`use "<CUSTOM_SUB_DOMAIN>.surge.sh" if you do not have a domain`)
				fmt.Println("to setup custom domain, see https://surge.sh/help/adding-a-custom-domain")
				return nil
			}

			domain := cCtx.Args().Get(1)

			if domain == "" {
				b, _ := os.ReadFile(path.Join(dir, "CNAME"))
				domain = strings.Trim(string(b), " ")
			}

			if domain == "" {
				fmt.Println(`Usage: surgecli upload <path_to_dir> <domain>`)
				fmt.Println("\nAs you have not specify a domain, generated a random domain for you")
				domain = fmt.Sprintf("%s.surge.sh", randomString(16))
				fmt.Printf("[%s] Accept it? (yes/no)  ", domain)
				confirmText1 := ""
				fmt.Scanf("%s\n", &confirmText1)
				if confirmText1 != "yes" {
					fmt.Println("\nAborted")
					return nil
				}
				cnameFilePath := path.Join(dir, "CNAME")
				fmt.Printf("Do you want write the domain to %s file?\n", cnameFilePath)
				fmt.Print("This will allow the cli to remember that domain  (yes/no)  ")
				confirmText2 := ""
				fmt.Scanf("%s\n", &confirmText2)
				if confirmText2 == "yes" {
					os.WriteFile(cnameFilePath, []byte(domain), 0666)
				}
			}

			fmt.Println()
			return c.surgesh.Upload(domain, dir, onUploadEvent)
		},
	}
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		bytes[i] = letters[r.Intn(len(letters))]
	}
	return string(bytes)
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
				fmt.Printf("%-7s <end>\n", p.Id)
			} else {
				fmt.Printf("%-7s %s \n", p.Id, p.File)
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
