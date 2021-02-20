// Copyright 2020-2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package mainthread_test

import (
	"fmt"

	"golang.design/x/mainthread"
)

func ExampleInit() {
	mainthread.Init(func() {
		// ... Do stuff ...
	})
	// Output:
}

func ExampleCall() {
	mainthread.Init(func() {
		mainthread.Call(func() {
			fmt.Println("from main thread")
		})
	})
	// Output: from main thread
}

func ExampleGo() {
	mainthread.Init(func() {
		done := make(chan string)
		mainthread.Go(func() {
			done <- "main thread"
		})
		fmt.Println("from", <-done)
	})
	// Output: from main thread
}
