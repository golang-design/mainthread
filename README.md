# mainthread

[![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/mainthread)](https://pkg.go.dev/golang.design/x/mainthread) [![Go Report Card](https://goreportcard.com/badge/golang.design/x/mainthread)](https://goreportcard.com/report/golang.design/x/mainthread)
![mainthread](https://github.com/golang-design/mainthread/workflows/mainthread/badge.svg?branch=master)

Package mainthread schedules function calls on the main thread in Go.

```
import "golang.design/x/mainthread"
```

## Quick Start

```go
package main

import "golang.design/x/mainthread"

func main() {
    mainthread.Init(func() {
        mainthread.Call(func() {
            // ... runs on the main thread ...
        })

        go func() {
            // ... runs concurrently ...
        }()

        mainthread.Call(func() {
            // ... runs on the main thread ...
        })
    })
}
```

## License

GNU GPLv3 &copy; 2020 The golang.design Initiative Authors