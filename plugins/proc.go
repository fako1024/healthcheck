package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/fako1024/healthcheck/errors"
)

// Proc denotes a Proc connection health check plugin
type Proc struct {
	name     string
	binaries []string
}

// NewProc instantiates a new Proc plugin
func NewProc() *Proc {
	return &Proc{
		name:     "proc",
		binaries: []string{},
	}
}

// RegisterFlags registers command line flags specific for the plugin
func (t *Proc) RegisterFlags() {
	pflag.StringArrayVar(&t.binaries, t.name+".binary", []string{}, "Process binary name (can be specified multiple times)")
}

// Run executes the Proc plugin
func (t *Proc) Run() (errs errors.Errors) {

	if len(t.binaries) == 0 {
		return
	}

	// Extract list of running processes
	processes, err := ioutil.ReadDir("/proc")
	if err != nil {
		return errors.Errors{
			fmt.Errorf("Error parsing system processes: %s", err),
		}
	}

	// Checkk all provided binaries
	for _, binary := range t.binaries {
		if err := t.checkBinary(processes, binary); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t *Proc) checkBinary(processes []os.FileInfo, expectedBinary string) error {
	for _, process := range processes {
		if process.IsDir() {
			if _, err := strconv.Atoi(process.Name()); err == nil {

				path := filepath.Join("/proc", process.Name(), "stat")
				if _, err := os.Stat(path); err == nil {

					statBytes, err := ioutil.ReadFile(path)
					if err != nil {
						return fmt.Errorf("Failed to read stat path %s: %s", path, err)
					}

					// Parse binary name from stat file
					statData := string(statBytes)
					binStart := strings.IndexRune(statData, '(') + 1
					binEnd := strings.IndexRune(statData[binStart:], ')')
					binary := statData[binStart : binStart+binEnd]

					if strings.Contains(binary, expectedBinary) {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("Binary %s not found to be running", expectedBinary)
}