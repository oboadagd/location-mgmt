// Package main implements a main routine that is the starting point of
// location-history-mgmt microservice.
package main

import (
	"github.com/oboadagd/location-mgmt/appconfig"
)

// main invokes method that start-up this microservice
func main() {
	appconfig.StartApp()
}
