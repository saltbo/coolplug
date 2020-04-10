package plugin

import (
	"context"
	"plugin"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

type Context struct {
	context.Context

	Router   *gin.RouterGroup
	Database *gorm.DB
	Cron     *cron.Cron
}

func NewContext(ctx context.Context, router *gin.RouterGroup, database *gorm.DB, cron *cron.Cron) *Context {
	return &Context{
		Context: ctx,

		Router:   router,
		Database: database,
		Cron:     cron,
	}
}

type Plugin interface {
	Name() string

	// It will be called when the plugin be installed
	Install() error

	// It will be called when the plugin be uninstalled
	Uninstall() error

	// It will be called when the system starts
	Run(c *Context) error

	// It will be called when the system shutdown
	Stop() error
}

func Load(filename string) (Plugin, error) {
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	s, err := p.Lookup("PluginConstructor")
	if err != nil {
		return nil, err
	}

	return s.(func() Plugin)(), nil
}
