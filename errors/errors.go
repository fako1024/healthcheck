package errors

import "strings"

const separator = "\n"

// Errors denotes a list of errors
type Errors []error

// Error returns a string representation of the error list
func (e Errors) Error() string {
	res := make([]string, len(e))
	for i := 0; i < len(e); i++ {
		res[i] = e.Error()
	}

	return strings.Join(res, separator)
}
