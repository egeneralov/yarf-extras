[![GoDoc](https://godoc.org/github.com/yarf-framework/extras/context/data?status.svg)](https://godoc.org/github.com/yarf-framework/extras/context/data)

## Context Data

Custom yarf.ContextData implementation package that works with strings. 
Includes a Yarf Middleware used to set the yarf.Data object automatically.

```go
import (
    "github.com/yarf-framework/yarf"
    "github.com/yarf-framework/extras/context/data"
    //...
)

func main() {
    y := yarf.New()
    y.Insert(new(data.SetStrData))
    
    // ...
}
```
