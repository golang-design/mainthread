// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a GNU GPLv3 license that can be found in the LICENSE file.

package mainthread // import "golang.design/x/mainthread"

import "runtime"

var funcQ = make(chan func(), runtime.GOMAXPROCS(0))

func init() {
	runtime.LockOSThread()
}

// Init initializes the functionality for running arbitrary subsequent
// functions on a main system thread.
//
// Init must be called in the main package.
func Init(main func()) {
	done := make(chan struct{})
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
	done := make(chan struct{})
	funcQ <- func() {
		defer func() {
			done <- struct{}{}
		}()
		f()
	}
	<-done
}
