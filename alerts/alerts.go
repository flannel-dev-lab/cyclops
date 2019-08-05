// alerts handles the alerting systems such as Sentry, etc
package alerts

// Interface to configure alerting system
type Alert interface {
	// Method to implement to capture error
	CaptureError(err error, message string)
	// Bootstrapping logic for alerts
	Bootstrap() error
}
