package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	p := PluginConstructor()
	assert.Equal(t, p.Name(), "example")
	assert.NoError(t, p.Install())
	assert.NoError(t, p.Uninstall())
}

func TestRandomInterval(t *testing.T) {
	RandomInterval(10, 30)
}
