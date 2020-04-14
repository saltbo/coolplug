package commands

import (
	"fmt"
	"net/url"

	"github.com/urfave/cli"

	"github.com/saltbo/coolplug/client"
	"github.com/saltbo/coolplug/cmd/client/flags"
)

func NewPluginCommands() cli.Commands {
	clientInitializer := func(c *cli.Context) client.Client {
		return client.NewHTTPClient(c.GlobalString(flags.FLAG_SERVER_HOST))
	}

	return cli.Commands{
		cli.Command{
			Name:   "list",
			Action: ActionRegister(clientInitializer, list),
		},
		cli.Command{
			Name:   "install",
			Action: ActionRegister(clientInitializer, install),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name",
					Required: true,
				},
				cli.StringFlag{
					Name:     "intro",
					Required: true,
				},
			},
		},
		cli.Command{
			Name:   "uninstall",
			Action: ActionRegister(clientInitializer, uninstall),
		},
	}
}

func list(c *ACtx) {
	mps, err := c.client.PluginList()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mps)
}

func install(c *ACtx) {
	for _, file := range c.cliCtx.Args() {
		pi := make(url.Values)
		pi.Set("name", c.cliCtx.String("name"))
		pi.Set("intro", c.cliCtx.String("intro"))
		err := c.client.PluginInstall(pi, file)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func uninstall(c *ACtx) {
	for _, id := range c.cliCtx.Args() {
		if err := c.client.PluginUninstall(id); err != nil {
			fmt.Println(err)
			return
		}
	}
}
