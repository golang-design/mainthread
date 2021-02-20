# mainthread [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/mainthread)](https://pkg.go.dev/golang.design/x/mainthread) ![mainthread](https://github.com/golang-design/mainthread/workflows/mainthread/badge.svg?branch=main) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/mainthread)

schedules function to run on the main thread

```go
import "golang.design/x/mainthread"
```

## Features

- Main thread scheduling
- Schedule functions without memory allocation

## API Usage

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

## When do you need this package?

Read this to learn more about the design purpose of this package:
https://golang.design/research/zero-alloc-call-sched/

## Who is using this package?

The initial purpose of building this package is to support writing
graphical applications in Go. To know projects that are using this
package, check our [wiki](https://github.com/golang-design/mainthread/wiki)
page.


## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).