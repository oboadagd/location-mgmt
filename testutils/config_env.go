// Package test implements utils for testutils of microservice location-mgmt.
package testutils

import (
	"strings"
)

// EvaluateErrConditions returns evaluation of values. Returns false if at least one evaluation is
// not met, returns true otherwise
func EvaluateErrConditions(errMsg string, values []string) bool {

	for _, value := range values {
		if !strings.Contains(strings.ToLower(errMsg), value) {
			return true
		}
	}

	return false
}
