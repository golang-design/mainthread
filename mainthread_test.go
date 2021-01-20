// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a GNU GPLv3 license that can be found in the LICENSE file.

package mainthread_test

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"golang.design/x/mainthread"
	"golang.org/x/sys/unix"
)

var initTid int

func init() {
	initTid = unix.Getpid()
}

func TestMain(m *testing.M) {
	mainthread.Init(func() {
		os.Exit(m.Run())
	})
}

func TestMainThread(t *testing.T) {
	var (
		nummain uint64
		numcall = 100000
	)

	wg := sync.WaitGroup{}
	for i := 0; i < numcall; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			mainthread.Call(func() {
				tid := unix.Gettid()
				if tid == initTid {
					return
				}
				t.Logf("call is not executed on the main thread, want %d, got %d", initTid, tid)
			})
		}()
		go func() {
			defer wg.Done()
			if unix.Gettid() == initTid {
				atomic.AddUint64(&nummain, 1)
			}
		}()
	}
	wg.Wait()

	if nummain == uint64(numcall) {
		t.Fatalf("all non main thread calls are executed on the main thread.")
	}
}

func BenchmarkCall(b *testing.B) {
	f := func() {}
	mainthread.Init(func() {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mainthread.Call(f)
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
