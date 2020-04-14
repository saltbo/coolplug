//  Copyright 2020 The Go Authors. All rights reserved.
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.

package commands

import (
	"github.com/urfave/cli"

	"github.com/saltbo/coolplug/client"
)

type Initializer func(c *cli.Context) client.Client

type ACtx struct {
	cliCtx *cli.Context
	client client.Client
}

type Action struct {
	initializer Initializer
	callback    func(c *ACtx)
}

func ActionRegister(initializer Initializer, callback func(c *ACtx)) func(c *cli.Context) {
	act := &Action{
		initializer: initializer,
		callback:    callback,
	}

	return func(c *cli.Context) {
		act.callback(&ACtx{
			cliCtx: c,
			client: act.initializer(c),
		})
	}
}
