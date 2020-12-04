# Handling Sessions

- Cyclops has the ability to handle sessions seamlessly
- Cyclops currenty supports `Redis` as the session store which implements the `Store` interface, which means that if Redis
is not your thing, you can always implement the methods in the `Store` interface and you can use whichever backend you like

### Initializing session
```go
package main

import (
	"github.com/flannel-dev-lab/cyclops/v2/cookie"
	"github.com/flannel-dev-lab/cyclops/v2/sessions"
)

func main() {
	store := sessions.RedisSessionStore{}
	err := store.New("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	cookieobj := cookie.CyclopsCookie{
		Secure:         false,
		HttpOnly:       true,
		StrictSameSite: false,
	}

	session := sessions.Session{
		Store:  &store,
		Cookie: cookieobj}
}
```

### Adding data to session
```
data := make(map[string]interface{})
data["Name"] = "hello"

session.Set(w, data, 10)
```
The above snippet says set the map `data` to session with 10 seconds expiry

### Get data from session
```
session.Get("Name")
```

### Delete data from session
```
session.Delete("Name")
```

### Delete all from backend session store
```
session.Reset()
```