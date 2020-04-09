package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/saltbo/coolplug/plugin"
)

type ExamplePlugin struct {
}

// PluginConstructor is required!!
func PluginConstructor() plugin.Plugin {
	return &ExamplePlugin{}
}

func (p *ExamplePlugin) Install(c *plugin.Context) error {
	c.Router.GET("/test", func(c *gin.Context) {

	})

	_, err := c.Cron.AddFunc("29 17 * * *", func() {
		for _, duration := range RandomInterval(10, 10) {
			<-time.After(duration)
			fmt.Println(duration)
		}
	})
	return err
}

func (p *ExamplePlugin) Uninstall() error {
	return nil
}

func RandomInterval(n, minSec int) []time.Duration {
	rand.Seed(time.Now().Unix())
	timeIntervals := make([]time.Duration, 0, n)
	for i := 0; i < n; i++ {
		timeIntervals = append(timeIntervals, time.Duration(rand.Intn(minSec)+minSec)*time.Second)
	}

	return timeIntervals
}
