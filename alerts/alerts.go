// Handles the alerting systems such as Sentry, etc
package alerts

import (
	"errors"
	sentry2 "github.com/flannel-dev-lab/cyclops/alerts/sentry"
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
		sentry := &sentry2.Sentry{
			DSN:         endpoint,
			Environment: environment,
			Enabled:     enabled,
		}
		return sentry, nil
	default:
		return nil, errors.New("undefined alerting system")
	}
}
