# Alerting

Alerting just allows you to use custom alerting providers such as sentry, etc. The only thing the 
developers need to do is to implement the alert interface

Example usage of alerting

```go
package main

import (
	"errors"
	"fmt"
	"github.com/flannel-dev-lab/cyclops/alerts"
	"github.com/flannel-dev-lab/cyclops/alerts/sentry"
	"net/http"
)

type Sentry struct {
	// URL for sentry
	DSN         string
	// Specifies if alert belongs to dev, stage, production environment
	Environment string
	// Trigger to see if alert is enabled or disabled
	Enabled     bool
}

func main() {
	sentry2 := sentry.Sentry{
		DSN: "",
		Environment: "dev",
		Enabled: true,
	}

	sentry2.Bootstrap()
	
	var alert alerts.Alert

	alert = sentry2
	alert.CaptureError(errors.New("test"), "test")
}


// CaptureError implements the alert interface to capture error
func (sentry Sentry) CaptureError(err error, message string) {
	if sentry.Enabled {
		raven.CaptureErrorAndWait(
			err,
			map[string]string{
				"env":     sentry.Environment,
				"message": message,
			})
	}
}

// Initializes the raven DSN
func (sentry Sentry) Bootstrap() error {
	return raven.SetDSN(sentry.DSN)
}

```