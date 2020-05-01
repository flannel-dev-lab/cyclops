# Alerting

Alerting allows you to use custom alerting providers such as sentry, etc. The only thing the 
developers need to do is to implement the alert interface

Example usage of alerting with sentry

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
	sentryAlerts := sentry.Sentry{
		DSN: "https://sentry.io/app1",
		Environment: "dev",
		Enabled: true,
	}

	sentryAlerts.Bootstrap()
	
	var alert alerts.Alert

	alert = sentryAlerts
	alert.CaptureError(errors.New("test"), "test")
}
```

As `Alert` is an interface{}, you can implement it to use your own custom alerting