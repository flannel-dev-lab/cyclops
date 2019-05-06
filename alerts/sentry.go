package alerts

import "github.com/getsentry/raven-go"

// Configuring Sentry alerting system
type Sentry struct {
	DSN         string
	Environment string
	Enabled     bool
}

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
