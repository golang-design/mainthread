# mainthread

[![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/mainthread)](https://pkg.go.dev/golang.design/x/mainthread) [![Go Report Card](https://goreportcard.com/badge/golang.design/x/mainthread)](https://goreportcard.com/report/golang.design/x/mainthread)
![mainthread](https://github.com/golang-design/mainthread/workflows/mainthread/badge.svg?branch=main)

Package mainthread schedules function to run on the main thread with zero allocation.

```go
import "golang.design/x/mainthread"
```

## Quick Start

```go
package main

import "golang.design/x/mainthread"

func main() {
    mainthread.Init(fn)
}

func fn() {
    mainthread.Call(func() {
        // ... runs on the main thread ...
    })

    // ... do what ever you want to do ...
}
```

## License

GNU GPLv3 &copy; 2020 The golang.design Initiative Authors