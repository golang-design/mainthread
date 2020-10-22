// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a GNU GPLv3 license that can be found in the LICENSE file.

package mainthread_test

import (
	"fmt"
	"testing"

	"golang.design/x/mainthread"
)

func TestMainthread(t *testing.T) {
	mainthread.Init(func() {
		mainthread.Call(func() {
			// FIXME: how to test this function really runs on the main thread?
		})
	})
}

func BenchmarkCall(b *testing.B) {
	f1 := func() {}
	f2 := func() {}

	mainthread.Init(func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				mainthread.Call(f1)
			} else {
				mainthread.Call(f2)
			}
		}
	})
}

func ExampleInit() {
	mainthread.Init(func() {
		mainthread.Call(func() {
			fmt.Println("from main thread")
		})
	})
	// Output: from main thread
}
