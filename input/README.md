# Handling user inputs

- Cyclops makes it easy to get query and form values.
- To get values from query parameters, do the following:
```
import "github.com/flannel-dev-lab/cyclops/v2/input"

func Login(w http.ResponseWriter, r *http.Request) {
    username := input.Query("username", r)
}
```
- The above code will capture query parameter `http://example.com?username="cyclops"`
- To get values from form data, do the following:
```
import "github.com/flannel-dev-lab/cyclops/v2/input"

func Login(w http.ResponseWriter, r *http.Request) {
    username := input.Form("username", r)
}
```

#  Handling File Uploads
Cyclops makes it seamless to handle file uploads using the method `input.FileContent(r *http.Request, key string)`. 
Here the `key` parameter is the name of the html tag to expect file content. 
All the user needs to do is to make sure that the form encoding supports multipart data. Consider the below example which describes how file handling works

```html
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>A static page</title>
    <link rel="stylesheet" href="/static/stylesheets/main.css" type="text/css">
</head>
<body>
<h1>Hello from a static page</h1>

<h1>Contact</h1>
<form action="/hello" method="POST" enctype="multipart/form-data">
    <input type="file" name="uploadfile" />
    <input type="submit">
</form>

</body>
</html>
```

```go
func FileHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := input.FileContent(r, "uploadfile")
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	return
}
```