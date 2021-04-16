package plugins

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"

	"github.com/fako1024/healthcheck/errors"
)

// TCP denotes a TCP connection health check plugin
type TCP struct {
	name      string
	endpoints []string
}

// NewTCP instantiates a new TCP plugin
func NewTCP() *TCP {
	return &TCP{
		name:      protoTCP,
		endpoints: []string{},
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *TCP) RegisterFlags() {
	pflag.StringArrayVar(&t.endpoints, t.name+".endpoint", []string{}, "TCP endpoint (can be specified multiple times)")
}

// Run executes the TCP plugin
func (t *TCP) Run() (errs errors.Errors) {

	// Checkk all provided endpoints
	for _, endpoint := range t.endpoints {
		if err := t.runEndpoint(endpoint); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t *TCP) runEndpoint(endpoint string) error {

	// Attempt to establish the TCP connection
	conn, err := net.Dial(protoTCP, endpoint)
	if err != nil {
		return fmt.Errorf("error establishing TCP connection to %s: %s", endpoint, err)
	}

	// Close the connection
	if err = conn.Close(); err != nil {
		return fmt.Errorf("error closing TCP connection to %s: %s", endpoint, err)
	}

	return nil
}
