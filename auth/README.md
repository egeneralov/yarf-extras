[![GoDoc](https://godoc.org/github.com/yarf-framework/extras/auth?status.svg)](https://godoc.org/github.com/yarf-framework/extras/auth)

# Auth package

Super-simple, yet secure, token-based auth package for Go's http.Request. 
Compatible with Yarf framework, includes Middleware ready to insert into your Yarf router. 


## Tokens

Tokens are generated from calculating the SHA512 hash from 256 bytes randomly generated and returned as a string encoded in UTF-8.
The result is a 128 characters long string all lower case representing the hash like: 

```
b6e184525010a39057878fb7d7eca73c39dde0ac8b2bcff26acd71034e5922d6b5a9e30923d5d35482df396e11e57df9adc085cdd47cd2b1095b1d2880f38d01
```


## Storage

The package uses an internal storage engine that consists in a in-memory (volatile) map.
Check the Storage interface to implement your own storage.


## Examples

### Create token:

```go
import (
    "github.com/yarf-framework/extras/auth"
    //...
)

func Login(username, password) string {
    // ...
    
    // Some user service login
    if user.Login(username, password) {
        // Create new token valid for 10 minutes and return it.
        return auth.NewToken(user.Id, 600) // 10 minutes token
    }
}
```


### Get, Validate and Refresh token

(This is what Auth middleware does)

```go
import (
    "github.com/yarf-framework/yarf"
    "github.com/yarf-framework/extras/auth"
    //...
)

func (sr *SomeResource) Get(c *yarf.Context) error {
    // Obtains Auth token from "Auth" header value.
    token := auth.GetToken(c.Request)
    
    // Validate token, return error if invalid
    data, err := auth.ValidateToken(token)
    if err != nil {
        return err
    }

    // Refresh token expiration when valid, if we want to.
    auth.RefreshToken(token)
    
    //...
}
```


### Custom storage

```go
import (
    "github.com/yarf-framework/extras/auth"
)

func SomeInitMethod() {
    // ...
    
    myStore := new(MyCustomStorageEngine)
    auth.RegisterStorage(myStore)
    
    // ...
}
```


### Set Yarf middleware

```go
import (
    "github.com/yarf-framework/yarf"
    "github.com/yarf-framework/extras/auth"
    //...
)

func main() {
    y := yarf.New()
    
    y.Insert(new(auth.Auth))
    
    //...
    
    y.Start(":80")
}
``` 


### Delete token

```go
import (
    "github.com/yarf-framework/yarf"
    "github.com/yarf-framework/extras/auth"
    //...
)

func (sr *SomeResource) Delete(c *yarf.Context) error {
    // Obtains Auth token from "Auth" header value.
    token := auth.GetToken(c.Request)
    
    // Delete token
    auth.DeleteToken(token)
    
    //...
}
```

