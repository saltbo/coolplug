package main

import (
	"log"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/urfave/cli"

	"github.com/saltbo/coolplug/core"
	"github.com/saltbo/coolplug/rest"
)

const (
	FLAG_DRIVER     = "driver"
	FLAG_DSN        = "dsn"
	FLAG_PLUGIN_DIR = "plugin-dir"
)

func main() {
	app := cli.NewApp()
	app.Compiled = time.Now()
	//app.Version = version.Long
	app.Name = "coolplug"
	app.Usage = "coolplug"
	app.Copyright = "(c) 2020 saltbo.cn"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  FLAG_DRIVER,
			Usage: "specify database driver, default: sqlite3",
			Value: "sqlite3",
		},
		cli.StringFlag{
			Name:  FLAG_DSN,
			Usage: "specify data source name, default: test.db",
			Value: "build/test.db",
		},
		cli.StringFlag{
			Name:  FLAG_PLUGIN_DIR,
			Usage: "specify plugin dir, default: build/plugins",
			Value: "build/plugins",
		},
	}
	app.Action = serverAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serverAction(c *cli.Context) {
	engine, err := core.New(c.String(FLAG_DRIVER), c.String(FLAG_DSN))
	if err != nil {
		log.Fatalln(err)
	}

	if err := engine.Boot(); err != nil {
		log.Fatalln(err)
	}

	rest.Boot(engine) // boot the rest server.
}
