package alerts

import (
	"errors"
)

// Interface to configure alerting system
type Alert interface {
	CaptureError(err error, message string)
}

// InitiateAlerting returns the new alerting object based on the system
// Supported systems are "sentry"
func InitiateAlerting(system, endpoint, environment string, enabled bool) (Alert, error) {
	switch system {
	case "sentry":
		sentry := &Sentry{
			DSN:         endpoint,
			Environment: environment,
			Enabled:     enabled,
		}
		return sentry, nil
	default:
		return nil, errors.New("undefined alerting system")
	}
}
