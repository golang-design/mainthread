// Copyright 2020-2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package mainthread offers facilities to schedule functions
// on the main thread. To use this package properly, one must
// call `mainthread.Init` from the main package. For example:
//
//	package main
//
//	import "golang.design/x/mainthread"
//
//	func main() { mainthread.Init(fn) }
//
//	// fn is the actual main function
//	func fn() {
//		// ... do stuff ...
//
//		// mainthread.Call returns when f1 returns. Note that if f1 blocks
//		// it will also block the execution of any subsequent calls on the
//		// main thread.
//		mainthread.Call(f1)
//
//		// ... do stuff ...
//
//
//		// mainthread.Go returns immediately and f2 is scheduled to be
//		// executed in the future.
//		mainthread.Go(f2)
//
//		// ... do stuff ...
//	}
//
//	func f1() { ... }
//	func f2() { ... }
//
// If the given function triggers a panic, and called via `mainthread.Call`,
// then the original panic value is propagated to the calling goroutine,
// preserving its type. One can capture that panic, when possible:
//
//	defer func() {
//		if r := recover(); r != nil {
//			println(r)
//		}
//	}()
//
//	mainthread.Call(func() { ... }) // if panic
//
// If the given function triggers a panic, and called via `mainthread.Go`,
// then the panic will be cached internally, until a call to the `Error()` method:
//
//	mainthread.Go(func() { ... }) // if panics
//
//	// ... do stuff ...
//
//	if err := mainthread.Error(); err != nil { // can be captured here.
//		println(err)
//	}
//
// Note that a panic happens before `mainthread.Error()` returning the
// panicked error. If one needs to guarantee `mainthread.Error()` indeed
// captured the panic, a dummy function can be used as synchornization:
//
//	mainthread.Go(func() { panic("die") })	// if panics
//	mainthread.Call(func() {}) 				// for execution synchronization
//	err := mainthread.Error()				// err must be non-nil
//
// It is possible to cache up to a maximum of 42 panicked errors.
// More errors are ignored.
package mainthread // import "golang.design/x/mainthread"

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {
	runtime.LockOSThread()
}

// Init initializes the functionality of running arbitrary subsequent
// functions be called on the main system thread.
//
// Init must be called in the main.main function.
func Init(main func()) {
	done := donePool.Get().(chan any)
	defer donePool.Put(done)

	go func() {
		defer func() {
			done <- nil
		}()
		main()
	}()

	for {
		select {
		case f := <-funcQ:
			dispatch(f)
		case <-done:
			return
		}
	}
}

// dispatch runs f on the calling (main) thread and reports any panic.
// For a Call, the recovered value is handed back verbatim through f.done
// so it can be re-panicked on the caller with its original type. For a Go,
// the panic is wrapped as an error and buffered for a later Error call.
func dispatch(f funcData) {
	defer func() {
		r := recover()
		if f.done != nil {
			f.done <- r
			return
		}
		if r != nil {
			select {
			case erroQ <- fmt.Errorf("%v", r):
			default:
			}
		}
	}()
	f.fn()
}

// Call calls f on the main thread and blocks until f finishes.
//
// If f panics, the original panic value is re-panicked on the calling
// goroutine, preserving its type.
func Call(f func()) {
	done := donePool.Get().(chan any)
	defer donePool.Put(done)

	funcQ <- funcData{fn: f, done: done}
	if r := <-done; r != nil {
		panic(r)
	}
}

// CallV calls f on the main thread, blocks until f finishes, and returns
// the value produced by f. It is the value-returning counterpart of Call.
//
// Unlike Call, CallV allocates: f is wrapped in a closure that captures the
// return value, so it is not suitable for allocation-sensitive hot paths.
// If f panics, the original panic value is re-panicked on the calling
// goroutine, preserving its type.
func CallV[T any](f func() T) (v T) {
	Call(func() { v = f() })
	return
}

// Go schedules f to be called on the main thread.
func Go(f func()) {
	funcQ <- funcData{fn: f}
}

// Error returns an error that is captured if there are any panics
// happened on the mainthread.
//
// It is possible to cache up to a maximum of 42 panicked errors.
// More errors are ignored.
func Error() error {
	select {
	case err := <-erroQ:
		return err
	default:
		return nil
	}
}

var (
	funcQ    = make(chan funcData, runtime.GOMAXPROCS(0))
	erroQ    = make(chan error, 42)
	donePool = sync.Pool{New: func() any {
		return make(chan any)
	}}
)

type funcData struct {
	fn   func()
	done chan any
}
