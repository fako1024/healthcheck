package plugins

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/pflag"

	"github.com/fako1024/healthcheck/errors"
)

// Ping denotes a Ping connection health check plugin
type Ping struct {
	name       string
	endpoints  []string
	count      int
	interval   time.Duration
	timeout    time.Duration
	privileged bool
}

// NewPing instantiates a new Ping plugin
func NewPing() *Ping {
	return &Ping{
		name:      "ping",
		endpoints: []string{},
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *Ping) RegisterFlags() {
	pflag.StringArrayVar(&t.endpoints, t.name+".endpoint", []string{}, "Ping endpoint (can be specified multiple times)")
	pflag.IntVar(&t.count, t.name+".count", 4, "Number of pings to send")
	pflag.DurationVar(&t.interval, t.name+".interval", 250*time.Millisecond, "Time interval between individual pings")
	pflag.DurationVar(&t.timeout, t.name+".timeout", 10*time.Second, "Timeout for the ping")
	pflag.BoolVar(&t.privileged, t.name+".privileged", false, "Use privileged mode to perform the ping")
}

// Run executes the Ping plugin
func (t *Ping) Run() (errs errors.Errors) {

	// Checkk all provided endpoints
	for _, endpoint := range t.endpoints {
		if err := t.runEndpoint(endpoint); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t *Ping) runEndpoint(endpoint string) error {

	// Attempt to instantiate the pinger
	pinger, err := ping.NewPinger(endpoint)
	if err != nil {
		return fmt.Errorf("error preparing ping to %s: %w", endpoint, err)
	}
	pinger.Count = t.count
	pinger.Interval = t.interval
	pinger.Timeout = t.timeout
	pinger.SetPrivileged(t.privileged)

	// Execute the ping
	if err = pinger.Run(); err != nil {
		return fmt.Errorf("error executing ping to %s: %w", endpoint, err)
	}

	// Check for packet loss
	stats := pinger.Statistics()
	if stats.PacketsRecv < stats.PacketsSent {
		return fmt.Errorf("encountered packet loss to %s: %d / %d probes received / sent", endpoint, stats.PacketsRecv, stats.PacketsSent)
	}

	return nil
}
