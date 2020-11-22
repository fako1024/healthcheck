package main

import (
	"fmt"
	"os"

	"github.com/fako1024/healthcheck/errors"
	"github.com/fako1024/healthcheck/plugins"
	"github.com/spf13/pflag"
)

// Manager dentoes a health check manager that takes care of all plugins
type Manager struct {
	plugins []plugins.Plugin
}

// NewManager instantiates a new plugin Manager
func NewManager() *Manager {
	return &Manager{
		plugins: plugins.AllPlugins,
	}
}

// RegisterFlags registers command line flags specific for all plugins
func (m *Manager) RegisterFlags() {
	for _, plugin := range m.plugins {
		plugin.RegisterFlags()
	}
	pflag.Parse()
}

// Run executes all plugins
func (m *Manager) Run() (errors errors.Errors) {
	for _, plugin := range m.plugins {
		if err := plugin.Run(); err != nil {
			errors = append(errors, err...)
		}
	}

	return errors
}

func main() {
	m := NewManager()

	m.RegisterFlags()

	if errs := m.Run(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}
