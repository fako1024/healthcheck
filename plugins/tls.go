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
		return fmt.Errorf("Invalid enpoint format: %s", endpoint)
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
	conn, err := tls.Dial("tcp", endpoint, &config)
	if err != nil {
		return fmt.Errorf("Error establishing TLS connection to %s: %s", endpoint, err)
	}

	state := conn.ConnectionState()
	if !state.HandshakeComplete || !state.NegotiatedProtocolIsMutual {
		conn.Close()
		return fmt.Errorf("Failed to complete TLS handshake: (complete: %v , mutual: %v)", state.HandshakeComplete, state.NegotiatedProtocolIsMutual)
	}

	// Close the connection
	if err = conn.Close(); err != nil {
		return fmt.Errorf("Error closing TLS connection to %s: %s", endpoint, err)
	}

	return nil
}
