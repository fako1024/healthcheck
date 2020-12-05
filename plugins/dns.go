package plugins

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/fako1024/healthcheck/errors"
)

// DNS denotes a DNS connection health check plugin
type DNS struct {
	name     string
	server   string
	queries  []string
	protocol string

	resolver *net.Resolver
}

// NewDNS instantiates a new DNS plugin
func NewDNS() *DNS {
	return &DNS{
		name:    "dns",
		queries: []string{},
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *DNS) RegisterFlags() {
	pflag.StringArrayVar(&t.queries, t.name+".query", []string{}, "DNS query (can be specified multiple times in the form <TYPE>;<REQUEST>;<EXPECTED RESULT>[;<PROTOCOL OVERRIDE>], e.g. AAAA;example.org;1.2.3.4,5.6.7.8[;udp6])")
	pflag.StringVar(&t.server, t.name+".server", "1.1.1.1:53", "DNS server to use for query")
	pflag.StringVar(&t.protocol, t.name+".protocol", "udp", `Protocol to use for the query ("tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only))`)
}

// Run executes the DNS plugin
func (t *DNS) Run() (errs errors.Errors) {

	t.resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, t.protocol, t.server)
		},
	}

	// Checkk all provided endpoints
	for _, query := range t.queries {
		if err := t.runEndpoint(query); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t *DNS) runEndpoint(query string) error {

	// Parse arguments
	args := strings.Split(query, ";")
	if len(args) < 3 {
		return fmt.Errorf("Invalid DNS query request, expected syntax <TYPE>;<REQUEST>;<EXPECTED RESULT>[;<PROTOCOL OVERRIDE>], got: %s", query)
	}

	// Override default resolver, if requested
	resolver := t.resolver
	if len(args) == 4 {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 10 * time.Second,
				}
				return d.DialContext(ctx, args[3], t.server)
			},
		}
	}

	// Perform query
	var entries []string
	switch args[0] {

	//Standard A / AAAA records
	case "A", "AAAA":
		network := "ip4"
		if args[0] == "AAAA" {
			network = "ip6"
		}
		ips, err := resolver.LookupIP(context.Background(), network, args[1])
		if err != nil {
			return err
		}
		for _, ip := range ips {
			entries = append(entries, ip.String())
		}

	// TXT records
	case "TXT":
		txt, err := resolver.LookupTXT(context.Background(), args[1])
		if err != nil {
			return err
		}
		entries = txt

	// Default: Unsupported
	default:
		return fmt.Errorf("Unsupported query type requested: %s", args[0])
	}

	// Trivial checks
	refEntries := strings.Split(args[2], ",")
	if len(entries) != len(refEntries) {
		return fmt.Errorf("Query result does not match expectation, want %v, have %v", refEntries, entries)
	}

	// Sort both reference + input slices for comparison
	sort.Slice(entries, func(i, j int) bool { return entries[i] < entries[j] })
	sort.Slice(refEntries, func(i, j int) bool { return refEntries[i] < refEntries[j] })
	if !stringsEqual(entries, refEntries) {
		return fmt.Errorf("Query result does not match expectation, want %v, have %v", refEntries, entries)
	}

	return nil
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
