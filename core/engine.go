package core

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/robfig/cron/v3"

	"github.com/saltbo/coolplug/plugin"
	"github.com/saltbo/coolplug/tools"
)

type Engine struct {
	Router   *gin.Engine
	Database *gorm.DB
	Cron     *cron.Cron

	// plugins store
	plugins []plugin.Plugin
}

func New(driver, dsn, pluginDir string) (*Engine, error) {
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		Router:   gin.Default(),
		Database: db,
		Cron:     cron.New(),
		plugins:  make([]plugin.Plugin, 0),
	}

	// load the plugins
	if err := tools.WalkAllFile(pluginDir, ".so", engine.Load); err != nil {
		return nil, err
	}

	return engine, nil
}

func (e *Engine) Load(filename string) error {
	p, err := plugin.Load(filename)
	if err != nil {
		return err
	}

	e.plugins = append(e.plugins, p)
	return nil
}


func (e *Engine) Boot() error {
	// install the plugins
	for _, p := range e.plugins {
		err := p.Install(plugin.NewContext(e.Router.Group("/external"), e.Database, e.Cron))
		if err != nil {
			return err
		}
	}

	// boot the cron daemon
	e.Cron.Start()
	return nil
}