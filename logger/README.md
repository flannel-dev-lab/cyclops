# Logger

Cyclops supports custom JSON logging if your application wants to add JSON logs

```go
package main

import (
    "github.com/flannel-dev-lab/cyclops/logger"
    "os"
)

func main() {
    log := logger.New(os.Stdout)
    
    log.Message(map[string]string{
    		"hello": "data",
    		"foo": "bar",
    	})
}
```

Output will be something like
```json
{"foo":"bar","hello":"data"}
```

Logger has following methods available:
```
- Message(message interface{})
- Fatal(message interface{})
- Panic(message interface{})
```

`logger.New()` takes in a writer interface, so in short, it can write to File, stderr, stdout etc which implement io.Writer
method
