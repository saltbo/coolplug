package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/saltbo/coolplug/plugin"
)

type ExamplePlugin struct {
	ctx *plugin.Context
}

// PluginConstructor is required!!
func PluginConstructor() plugin.Plugin {
	return &ExamplePlugin{}
}

func (p *ExamplePlugin) Name() string {
	return "example"
}

func (p *ExamplePlugin) Install() error {
	//p.ctx.Database.CreateTable()
	return nil
}

func (p *ExamplePlugin) Uninstall() error {
	//p.ctx.Database.DropTable()
	return nil
}

func (p *ExamplePlugin) Run(c *plugin.Context) error {
	//time.Sleep(time.Second * 2)
	//select {
	//case <-c.Done():
	//	return c.Err()
	//default:
	//
	//}

	// register a router.
	c.Router.GET("/test", func(c *gin.Context) {
		fmt.Println(c.Request)
	})

	// register a cron task
	_, err := c.Cron.AddFunc("29 17 * * *", func() {
		for _, duration := range RandomInterval(10, 10) {
			<-time.After(duration)
			fmt.Println(duration)
		}
	})

	p.ctx = c
	return err
}

func (p *ExamplePlugin) Stop() error {
	//p.ctx.Cron.Remove()
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
