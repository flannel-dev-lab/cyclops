// Package alerts handles the alerting systems such as Sentry, etc
package alerts

// Alert is an interface to configure alerting system. Cyclops provides sentry alerting by default and implements the
// following methods
type Alert interface {
	// CaptureError captures error and a message and sends an alert
	CaptureError(err error, message string)
	// Bootstrap acts like a constructor to setup alerting parameters such as environment, alerting endpoint, etc.
	Bootstrap() error
}
