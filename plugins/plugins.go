package plugins

import (
	"github.com/sirupsen/logrus"
)

// PluginContext contains all required information for a plugin to run
type PluginContext struct {
	OutputPath string
	Config     map[string]string
	Logger     *logrus.Entry
	Finish     chan bool
}

// Plugin defines the interface required for all evacuate plugins
type Plugin interface {
	ShouldRun(PluginContext)
	Run(PluginContext)
}

// List TODO
type List map[string]Plugin

var _plugins List

// RegisterPlugin TODO
func RegisterPlugin(name string, plugin Plugin) {
	if _plugins == nil {
		_plugins = make(List)
	}

	_plugins[name] = plugin
}

// GetPlugins returns a list of registered plugins
func GetPlugins() List {
	return _plugins
}
