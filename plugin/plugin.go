package plugin

import (
	"plugin"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

type Context struct {
	Router   *gin.RouterGroup
	Database *gorm.DB
	Cron     *cron.Cron
}

func NewContext(router *gin.RouterGroup, database *gorm.DB, cron *cron.Cron) *Context {
	return &Context{
		Router:   router,
		Database: database,
		Cron:     cron,
	}
}

type Plugin interface {
	Install(c *Context) error
	Uninstall() error
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
