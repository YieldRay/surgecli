package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	surgeUtils "github.com/yieldray/surgecli/utils"
)

func (c *privateSurgeCLI) ProxyCommand() *cli.Command {
	return &cli.Command{
		Name:      "set-api",
		Usage:     "Set api host, the official host is https://surge.surge.sh",
		ArgsUsage: "<host>",
		Action: func(cCtx *cli.Context) error {
			api := cCtx.Args().First()
			if len(api) == 0 {
				fmt.Printf("e.g.  %s set-api https://surge.surge.sh\n", os.Args[0])
			}
			if err := surgeUtils.ConfSetApi(api); err != nil {
				return err
			} else {
				fmt.Println("Your api host has been reset to the official host")
				return nil
			}
		},
	}
}
