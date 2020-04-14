package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/saltbo/coolplug/cmd/client/commands"
	"github.com/saltbo/coolplug/cmd/client/flags"
)

func initCommands() cli.Commands {
	return cli.Commands{
		cli.Command{
			Name:        "plugin",
			Subcommands: commands.NewPluginCommands(),
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Compiled = time.Now()
	//app.Version = version.Long
	app.Name = "coolctl"
	app.Usage = "coolctl"
	app.Copyright = "(c) 2020 saltbo.cn"
	app.Flags = flags.New()
	app.Commands = initCommands()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
