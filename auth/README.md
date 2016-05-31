[![GoDoc](https://godoc.org/github.com/yarf-framework/extras/auth?status.svg)](https://godoc.org/github.com/yarf-framework/extras/auth)

# Auth package

Super-simple, yet secure, token-based auth package for Go's http.Request. 
Compatible with Yarf framework, includes Middleware ready to insert into your Yarf router. 


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


## Set Yarf middleware

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

