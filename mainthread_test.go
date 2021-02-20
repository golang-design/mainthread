// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux
// +build linux

package mainthread_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.design/x/mainthread"
	"golang.org/x/sys/unix"
)

var initTid int

func init() {
	initTid = unix.Getpid()
}

func TestMain(m *testing.M) {
	mainthread.Init(func() { os.Exit(m.Run()) })
}

// TestMainThread is not designed to be executed on the main thread.
// This test tests the a call from this function that is invoked by
// mainthread.Call is either executed on the main thread or not.
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
				// Code inside this function is expecting to be executed
				// on the mainthread, this means the thread id should be
				// euqal to the initial process id.
				tid := unix.Gettid()
				if tid == initTid {
					return
				}
				t.Errorf("call is not executed on the main thread, want %d, got %d", initTid, tid)
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

func TestGo(t *testing.T) {
	done := make(chan struct{})
	mainthread.Go(func() {
		time.Sleep(time.Second)
		done <- struct{}{}
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	select {
	case <-ctx.Done():
	case <-done:
		t.Fatalf("mainthread.Go is not executing in parallel")
	}

	ctxx, cancell := context.WithTimeout(context.Background(), time.Second)
	defer cancell()
	select {
	case <-ctxx.Done():
		t.Fatalf("mainthread.Go never schedules the function")
	case <-done:
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
