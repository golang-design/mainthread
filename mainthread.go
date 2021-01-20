// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a GNU GPLv3 license that can be found in the LICENSE file.

package mainthread // import "golang.design/x/mainthread"

import (
	"runtime"
	"sync"
)

var funcQ = make(chan func(), runtime.GOMAXPROCS(0))

func init() {
	runtime.LockOSThread()
}

// Init initializes the functionality for running arbitrary subsequent
// functions on a main system thread.
//
// Init must be called in the main package.
func Init(main func()) {
	done := donePool.Get().(chan struct{})
	defer donePool.Put(done)

	go func() {
		defer func() {
			done <- struct{}{}
		}()
		main()
	}()

	for {
		select {
		case f := <-funcQ:
			f()
		case <-done:
			return
		}
	}
}

// Call calls f on the main thread and blocks until f finishes.
func Call(f func()) {
	done := donePool.Get().(chan struct{})
	defer donePool.Put(done)

	funcQ <- func() {
		defer func() {
			done <- struct{}{}
		}()
		f()
	}
	<-done
}

var donePool = sync.Pool{
	New: func() interface{} { return make(chan struct{}) },
}
