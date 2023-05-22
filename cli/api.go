package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	surgeUtils "github.com/yieldray/surgecli/utils"
)

func (c *privateSurgeCLI) ProxyCommand() *cli.Command {
	return &cli.Command{
		Name:      "api",
		Usage:     "Print or set api host, the official host is https://surge.surge.sh",
		ArgsUsage: "<host>",
		Action: func(cCtx *cli.Context) error {
			api := cCtx.Args().First()
			prevApi := surgeUtils.ConfGetApi()
			if len(api) == 0 {
				fmt.Printf("Usage: %s api https://surge.surge.sh\n", os.Args[0])
				fmt.Printf("Current api host is %s\n", prevApi)
				return nil
			}

			if err := surgeUtils.ConfSetApi(api); err != nil {
				return err
			} else {
				fmt.Printf("The previous api host is %s\n", prevApi)
				fmt.Printf("Your api host has been set to %s\n", api)
				return nil
			}
		},
	}
}
