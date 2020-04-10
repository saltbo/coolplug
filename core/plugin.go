package core

import (
	"fmt"
	"sync"

	"github.com/saltbo/coolplug/plugin"
)

type PluginInstaller struct {
	store sync.Map // name => Plugin
}

func NewPluginInstaller() *PluginInstaller {
	return &PluginInstaller{}
}

func (p *PluginInstaller) Load(filename string) error {
	instance, err := plugin.Load(filename)
	if err != nil {
		return err
	}

	p.store.Store(filename, instance)
	return nil
}

func (p *PluginInstaller) Install(filename string) (plugin.Plugin, error) {
	_, ok := p.store.Load(filename)
	if ok {
		return nil, fmt.Errorf("plugin [%s] already exist, please uninstall first.", filename)
	}

	instance, err := plugin.Load(filename)
	if err != nil {
		return nil, err
	}

	if err := instance.Install(); err != nil {
		return nil, err
	}

	p.store.Store(filename, instance)
	return instance, nil
}

func (p *PluginInstaller) Uninstall(filename string) error {
	v, ok := p.store.Load(filename)
	if !ok {
		return fmt.Errorf("plugin [%s] not exist", filename)
	}

	instance := v.(plugin.Plugin)
	if err := instance.Stop(); err != nil {
		return err
	}

	if err := instance.Uninstall(); err != nil {
		return err
	}

	p.store.Delete(filename)
	return nil
}

func (p *PluginInstaller) Get(name string) (plugin.Plugin, error) {
	v, ok := p.store.Load(name)
	if !ok {
		return nil, fmt.Errorf("plugin [%s] not exist", name)
	}

	return v.(plugin.Plugin), nil
}

func (p *PluginInstaller) Range(f func(name string, plugin plugin.Plugin) bool) {
	p.store.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(plugin.Plugin))
	})
}
