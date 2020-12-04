package plugins

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/fako1024/healthcheck/errors"

	// Import SQL drivers
	_ "github.com/go-sql-driver/mysql"
)

// SQL denotes a SQL connection health check plugin
type SQL struct {
	name   string
	driver string
	host   string
	port   string
}

// NewSQL instantiates a new SQL plugin
func NewSQL() *SQL {
	return &SQL{
		name: "sql",
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *SQL) RegisterFlags() {
	pflag.StringVar(&t.driver, t.name+".driver", "mysql", "SQL driver")
	pflag.StringVar(&t.host, t.name+".host", "", "SQL host")
	pflag.StringVar(&t.port, t.name+".port", "3306", "SQL port")
}

// Run executes the SQL plugin
func (t *SQL) Run() (errs errors.Errors) {

	if t.host == "" {
		return
	}

	db, err := sql.Open(t.driver, "tcp("+t.host+":"+t.port+")/")
	if err != nil {
		return errors.Errors{
			fmt.Errorf("Error establishing SQL connection to %s: %s", "tcp("+t.host+":"+t.port+")/", err),
		}
	}
	defer db.Close()

	_, err = db.Exec("DO 1")
	if err != nil && !strings.Contains(err.Error(), ": Access denied for user") {
		return errors.Errors{
			fmt.Errorf("Unexpected error performing base SQL query: %s", err),
		}
	}

	return
}
