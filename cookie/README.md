# Handling Cookies

- Cookies in cyclops are described using the structure `CyclopsCookie`
- The fields in the structure are as follows:

| Attribute   | Optional  |
|-------------|:---------:|
| Name        | No |
| Value       | No |
| Path        | Yes |
| Domain      | No |
| Secure      | Yes |
| HttpOnly    | Yes |
| SameSite    | Yes |
| Expires | Yes |
| MaxAge | Yes |

## Creating a Cookie:
```
import "github.com/flannel-dev-lab/cyclops/cookie"

func Login(w http.ResponseWriter, r *http.Request) {
    cookieObj := cookie.CyclopsCookie{
            Name:           "test",
            Value:          "test",
            Secure:         false,
            HttpOnly:       true,
            StrictSameSite: false,
        }
    
    cookieObj.SetCookie(w)
}
```
Once the cookie is created, we use the method `SetCookie` which takes in a `http.ResponseWriter` to write the cookie

## Reading a cookie:
```
import "github.com/flannel-dev-lab/cyclops/cookie"

func Login(w http.ResponseWriter, r *http.Request) {
    cookieObj := cookie.CyclopsCookie{}

    fmt.Println(cookieObj.GetCookie(r, "test"))
}

```
To read a cookie, we create an empty cookie object and call the `GetCookie` method which takes in `*http.Request`, the 
`Name` of the cookie and returns a `*http.Cookie`

## Reading all cookies
```
import "github.com/flannel-dev-lab/cyclops/cookie"

func Login(w http.ResponseWriter, r *http.Request) {
    cookieObj := cookie.CyclopsCookie{}

    fmt.Println(cookieObj.GetAll(r, "test"))
}
```
To read a cookie, we create an empty cookie object and call the `GetAll` method which takes in `*http.Request` and 
returns a array  of `*http.Cookie`

## Deleting a Cookie
Deletes a cookie by setting max-age to 0
```
import "github.com/flannel-dev-lab/cyclops/cookie"

func Login(w http.ResponseWriter, r *http.Request) {
    cyclopsCookie := cookie.CyclopsCookie{}

    cookie, _ := cookieObj.GetCookie(r, "test")
    
    cyclopsCookie.Delete(w, cookie)
}
```
