package plugins

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"

	"github.com/fako1024/healthcheck/errors"
)

// HTTP denotes a HTTP connection health check plugin
type HTTP struct {
	name               string
	uri                string
	method             string
	expectedStatusCode int
	timeout            time.Duration
}

// NewHTTP instantiates a new HTTP plugin
func NewHTTP() *HTTP {
	return &HTTP{
		name: "http",
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *HTTP) RegisterFlags() {
	pflag.StringVar(&t.uri, t.name+".uri", "", "HTTP URI")
	pflag.StringVar(&t.method, t.name+".method", "GET", "HTTP method")
	pflag.IntVar(&t.expectedStatusCode, t.name+".expectedStatusCode", http.StatusOK, "Expected HTTP status code")
	pflag.DurationVar(&t.timeout, t.name+".timeout", 5*time.Second, "HTTP request timeout")
}

// Run executes the HTTP plugin
func (t *HTTP) Run() (errs errors.Errors) {

	if t.uri == "" {
		return
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(t.uri)
	req.Header.SetMethod(t.method)

	// Execute request
	if err := fasthttp.DoTimeout(req, resp, t.timeout); err != nil {
		return errors.Errors{
			fmt.Errorf("error performing request to %s: %w", t.uri, err),
		}
	}

	if resp.StatusCode() != t.expectedStatusCode {
		return errors.Errors{
			fmt.Errorf("unexpected HTTP status code in call to %s, want %d, have %d", t.uri, t.expectedStatusCode, resp.StatusCode()),
		}
	}

	return
}
