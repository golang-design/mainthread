// Copyright 2020-2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package mainthread_test

import (
	"context"
	"os"
	"testing"
	"time"

	"golang.design/x/mainthread"
)

// TestMain runs the entire test suite inside the main thread event loop so
// that mainthread.Call/Go have a loop to drain on every platform.
func TestMain(m *testing.M) {
	mainthread.Init(func() { os.Exit(m.Run()) })
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

	ctxx, cancell := context.WithTimeout(context.Background(), time.Second*2)
	defer cancell()
	select {
	case <-ctxx.Done():
		t.Fatalf("mainthread.Go never schedules the function")
	case <-done:
	}
}

func TestCallV(t *testing.T) {
	got := mainthread.CallV(func() int { return 42 })
	if got != 42 {
		t.Fatalf("CallV returned %d, want 42", got)
	}
}

func TestPanickedFuncCall(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Fatalf("expected to panic, but actually not")
	}()

	mainthread.Call(func() {
		panic("die")
	})
}

// TestPanickedFuncCallValue verifies that Call re-panics the original value,
// preserving its type rather than flattening it to a string.
func TestPanickedFuncCallValue(t *testing.T) {
	type myErr struct{ code int }

	defer func() {
		r := recover()
		got, ok := r.(myErr)
		if !ok {
			t.Fatalf("expected panic value of type myErr, got %T (%v)", r, r)
		}
		if got.code != 7 {
			t.Fatalf("expected code 7, got %d", got.code)
		}
	}()

	mainthread.Call(func() {
		panic(myErr{code: 7})
	})
}

func TestPanickedFuncCallV(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Fatalf("expected to panic, but actually not")
	}()

	mainthread.CallV(func() int {
		panic("die")
	})
}

func TestPanickedFuncGo(t *testing.T) {
	defer func() {
		if err := mainthread.Error(); err != nil {
			return
		}
		t.Fatalf("expected to panic, but actually not")
	}()

	mainthread.Go(func() { panic("die") })
	mainthread.Call(func() {}) // for sync
}
