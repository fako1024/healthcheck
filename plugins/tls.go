package plugins

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/fako1024/healthcheck/errors"
)

// TLS denotes a TLS connection health check plugin
type TLS struct {
	name      string
	endpoints []string
}

// NewTLS instantiates a new TLS plugin
func NewTLS() *TLS {
	return &TLS{
		name:      "tls",
		endpoints: []string{},
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *TLS) RegisterFlags() {
	pflag.StringArrayVar(&t.endpoints, t.name+".endpoint", []string{}, "TLS endpoint (can be specified multiple times)")
}

// Run executes the TLS plugin
func (t *TLS) Run() (errs errors.Errors) {

	// Checkk all provided endpoints
	for _, endpoint := range t.endpoints {
		if err := t.runEndpoint(endpoint); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t *TLS) runEndpoint(endpoint string) error {

	endpointParams := strings.Split(endpoint, ":")
	if len(endpointParams) < 2 || len(endpointParams) > 3 {
		return fmt.Errorf("invalid enpoint format: %s", endpoint)
	}

	// Initiate TLS parameters and override hostname if required
	config := tls.Config{
		InsecureSkipVerify: false,
	}
	if len(endpointParams) == 3 {
		config.ServerName = endpointParams[0]
		endpoint = endpointParams[1] + ":" + endpointParams[2]
	}

	// Attempt to establish the TLS connection
	conn, err := tls.Dial(protoTCP, endpoint, &config)
	if err != nil {
		return fmt.Errorf("error establishing TLS connection to %s: %s", endpoint, err)
	}

	state := conn.ConnectionState()
	if !state.HandshakeComplete {
		err = conn.Close()
		return fmt.Errorf("failed to complete TLS handshake: (complete: %v), close err: %s", state.HandshakeComplete, err)
	}

	// Close the connection
	if err = conn.Close(); err != nil {
		return fmt.Errorf("error closing TLS connection to %s: %s", endpoint, err)
	}

	return nil
}
