// Package sentry implements the Alert interface for Sentry alerting systems
package alerts

import "github.com/getsentry/raven-go"

// Sentry struct to hold sentry variables
type Sentry struct {
	// DSN is the URL for sentry
	DSN string
	// Environment specifies if alert belongs to dev, stage, production environment
	Environment string
	// Enabled is a trigger for an alert to be enabled or disabled
	Enabled bool
}

// CaptureError implements the alert interface to capture error and send a message if alerting is enabled
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

// Bootstrap initializes the raven DSN by setting the environment
func (sentry Sentry) Bootstrap() error {
	raven.SetEnvironment(sentry.Environment)
	return raven.SetDSN(sentry.DSN)
}
