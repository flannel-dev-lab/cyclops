// Sentry alerting system
package sentry

import "github.com/getsentry/raven-go"

// Struct to hold sentry environment variables
type Sentry struct {
	// URL for sentry
	DSN         string
	// Specifies if alert belongs to dev, stage, production environment
	Environment string
	// Trigger to see if alert is enabled or disabled
	Enabled     bool
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
