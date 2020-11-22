package plugins

import "github.com/fako1024/healthcheck/errors"

// AllPlugins denotes a list of all available plugins
var AllPlugins = []Plugin{
	NewTCP(),
	NewSSH(),
	NewHTTP(),
	NewSQL(),
	NewProc(),
	NewTLS(),
}

// Plugin denotes a generic health check plugin
type Plugin interface {

	// RegisterFlags registers command line flags specific for the plugin
	RegisterFlags()

	// Run executes the plugin
	Run() errors.Errors
}
