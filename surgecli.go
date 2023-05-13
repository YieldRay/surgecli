package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/surge"
	"github.com/yieldray/surgecli/types"
)

const DEBUG = false

type Rt struct {
}

func (rt Rt) RoundTrip(req *http.Request) (*http.Response, error) {
	if DEBUG {
		fmt.Println(req)
	}
	res, err := http.DefaultClient.Do(req)
	if DEBUG {
		fmt.Println(res)
	}
	return res, err
}

func main() {
	surgesh := surge.New()
	surgesh.SetHTTPClient(&http.Client{Transport: &Rt{}})

	app := &cli.App{
		Name:  "surgecli",
		Usage: "thrid party surge.sh cli",
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login (or create new account) to surge.sh",
				Action: func(cCtx *cli.Context) error {
					username := cCtx.Args().Get(0)
					password := cCtx.Args().Get(1)
					email, err := surgesh.Login(username, password)
					if err == nil {
						fmt.Printf("Login Success as %s\n", email)
					}
					return err
				},
			},
			{
				Name:  "logout",
				Usage: "logout to surge.sh",
				Action: func(cCtx *cli.Context) error {
					if email := surgesh.Whoami(); email == "" {
						fmt.Println("<YOU ARE NOT LOGGED IN>")
						return nil
					}
					email, err := surgesh.Logout()
					if err == nil {
						fmt.Printf("Logout Success as %s\n", email)
					}
					return err
				},
			},
			{
				Name:  "account",
				Usage: "show account information",
				Action: func(cCtx *cli.Context) error {
					if email := surgesh.Whoami(); email == "" {
						fmt.Println("<YOU ARE NOT LOGGED IN>")
						return nil
					}
					ac, err := surgesh.Account()
					if err != nil {
						return err
					}
					fmt.Println(ac)

					return nil
				},
			},
			{
				Name:  "whoami",
				Usage: "show my email",
				Action: func(cCtx *cli.Context) error {
					if email := surgesh.Whoami(); email == "" {
						fmt.Println("<YOU ARE NOT LOGGED IN>")
					} else {
						fmt.Println(email)
					}
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list my sites",
				Action: func(cCtx *cli.Context) error {
					list, err := surgesh.List()
					if err != nil {
						return err
					}
					for _, site := range list {
						fmt.Println(site)
					}
					return nil
				},
			},
			{
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
						return surgesh.Upload(domain, dir, onUploadEvent)
					}

				},
			},
			{
				Name:  "teardown",
				Usage: "delete site from surge.sh",
				Action: func(cCtx *cli.Context) error {
					domain := cCtx.Args().Get(0)

					if domain == "" {
						fmt.Println("please specify a domain to teardown")
						return nil
					}

					teardown, err := surgesh.Teardown(domain)
					if err != nil {
						return err
					}

					fmt.Println(teardown)

					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "API_HOST",
				Value: "https://surge.surge.sh",
				Usage: "configure the API host",
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
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
