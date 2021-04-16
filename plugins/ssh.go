package plugins

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"

	"github.com/fako1024/healthcheck/errors"
)

// SSH denotes a SSH connection health check plugin
type SSH struct {
	name      string
	endpoints []string
}

// NewSSH instantiates a new SSH plugin
func NewSSH() *SSH {
	return &SSH{
		name:      "ssh",
		endpoints: []string{},
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *SSH) RegisterFlags() {
	pflag.StringArrayVar(&t.endpoints, t.name+".endpoint", []string{}, "SSH endpoint (can be specified multiple times)")
}

// Run executes the SSH plugin
func (t *SSH) Run() (errs errors.Errors) {

	// Checkk all provided endpoints
	for _, endpoint := range t.endpoints {
		if err := t.runEndpoint(endpoint); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t *SSH) runEndpoint(endpoint string) error {

	sshConfig := &ssh.ClientConfig{
		User: "test",
		Auth: []ssh.AuthMethod{
			ssh.Password("test"),
		},
		/* #nosec */
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	_, err := ssh.Dial(protoTCP, endpoint, sshConfig)
	if err != nil && !strings.Contains(err.Error(), "ssh: unable to authenticate, attempted methods") {
		return fmt.Errorf("error establishing SSH connection to endpoint %s: %s", endpoint, err)
	}

	return nil
}
