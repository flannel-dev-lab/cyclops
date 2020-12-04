# Logger

Cyclops provides lightweight logger. It simple to use and is based on contexts

```go
package main

import (
	"net/http"
	"github.com/flannel-dev-lab/cyclops/v2/logger"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = logger.AddKey(ctx, "hello", "world")

	logger.Info(ctx, "info")
}

```

The above code will printout

```json
{
  "INFO": "info",
  "hello": "world"
}
```

You can add any number of keys to the log you want to track. There are other log levels available as well
`Debug`, `Warn`, `Error` that can be tuned as per your use case
