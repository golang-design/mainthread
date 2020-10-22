// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a GNU GPLv3 license that can be found in the LICENSE file.

package mainthread // import "golang.design/x/mainthread"

import "runtime"

// FIXME: can we do scheduling with zero overhead?
var fqueue chan func()

func init() {
	runtime.LockOSThread()

	// FIXME: what else can we do about queue size?
	fqueue = make(chan func(), runtime.GOMAXPROCS(0))
}

// Init initializes the functionality for running arbitrary subsequent
// functions on a main system thread.
//
// Init must be called in the main package.
func Init(run func()) {
	done := make(chan struct{})
	go func() {
		defer func() {
			// FIXME: do something about panicked f.
			recover()

			done <- struct{}{}
		}()
		run()
	}()

	for {
		select {
		case f := <-fqueue:
			f()
		case <-done:
			return
		}
	}
}

// Call calls f on the main thread and blocks until f finishes.
func Call(f func()) {
	done := make(chan struct{})
	fqueue <- func() {
		defer func() {
			// FIXME: do something about panicked f.
			recover()

			done <- struct{}{}
		}()
		f()
	}
	<-done
}
