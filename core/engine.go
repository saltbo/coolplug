package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"

	"github.com/saltbo/coolplug/model"
	"github.com/saltbo/coolplug/plugin"
)

type Engine struct {
	Router   *gin.Engine
	Database *gorm.DB
	Cron     *cron.Cron

	// plugin installer
	installer *PluginInstaller

	// plugin instance channel
	instanceCh chan plugin.Plugin
}

func New(driver, dsn string) (*Engine, error) {
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		Router:     gin.Default(),
		Database:   db,
		Cron:       cron.New(),
		installer:  NewPluginInstaller(),
		instanceCh: make(chan plugin.Plugin),
	}

	return engine, nil
}

func (e *Engine) PluginInstall(mp *model.Plugin) error {
	if err := e.Database.First(&model.Plugin{Name: mp.Name}).Error; err == nil {
		return fmt.Errorf("plugin [%s] already installed", mp.Name)
	}

	instance, err := e.installer.Install(mp.Filename)
	if err != nil {
		return err
	}

	e.instanceCh <- instance
	return e.Database.Create(mp).Error
}

func (e *Engine) PluginUninstall(name string) error {
	if err := e.installer.Uninstall(name); err != nil {
		return err
	}

	return e.Database.Delete(&model.Plugin{Name: name}).Error
}

func (e *Engine) pluginCtx(ctx context.Context) *plugin.Context {
	router := e.Router.Group("/external")
	return plugin.NewContext(ctx, router, e.Database, e.Cron)
}

func (e *Engine) runPlugin(instance plugin.Plugin) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan int, 1)
	go func() {
		if err := instance.Run(e.pluginCtx(ctx)); err != nil {
			log.Printf("plugin %s run error: %s\n", instance.Name(), err)
		}
		done <- 1
	}()

	select {
	case <-ctx.Done():
		log.Printf("plugin [%s] run timeout\n", instance.Name())
		return
	case <-done:
	}
}

func (e *Engine) Boot() error {
	// auto migrate the database model
	e.Database.AutoMigrate(&model.Plugin{})

	// load the all installed plugins
	plugins := make([]model.Plugin, 0)
	if err := e.Database.Find(&plugins).Error; err != nil {
		return err
	}
	for _, mp := range plugins {
		if err := e.installer.Load(mp.Filename); err != nil {
			return err
		}
	}

	// run the installed plugins
	go func() {
		for instance := range e.instanceCh {
			e.runPlugin(instance)
		}
	}()

	e.installer.Range(func(name string, plugin plugin.Plugin) bool {
		e.instanceCh <- plugin
		return true
	})

	// boot the cron daemon
	e.Cron.Start()
	return nil
}
