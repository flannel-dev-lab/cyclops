# Handling user inputs

- Cyclops makes it easy to get query and form values.
- To get values from query parameters, do the following:
```
import "github.com/flannel-dev-lab/cyclops/input"

func Login(w http.ResponseWriter, r *http.Request) {
    username := input.Query("username", r)
}
```
- The above code will capture query parameter `http://example.com?username="cyclops"`
- To get values from form data, do the following:
```
import "github.com/flannel-dev-lab/cyclops/input"

func Login(w http.ResponseWriter, r *http.Request) {
    username := input.Form("username", r)
}
```
