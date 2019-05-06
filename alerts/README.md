# Alerting

Example usage of alerting

```go
func alerting() {
	alert, _ := alerts.InitiateAlerting("sentry", "http://sentry.io")
	alert.CaptureError(errors.New("hello"), "hello", "dev", true)
}
```