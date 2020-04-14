package flags

import (
	"github.com/urfave/cli"
)

const (
	FLAG_SERVER_HOST = "server_host"
)

func New() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  FLAG_SERVER_HOST,
			Usage: "specify server host",
			Value: "http://localhost:8080",
		},
	}
}
