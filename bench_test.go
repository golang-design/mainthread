// Copyright 2020-2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package mainthread_test

import (
	"testing"

	"golang.design/x/mainthread"
)

// BenchmarkCall relies on the main thread event loop started by TestMain,
// so it measures a single Call round-trip to the real main OS thread rather
// than a nested loop. The headline metric is allocations per op (zero); the
// ns/op reflects the cost of waking the locked main thread and is inherently
// platform-dependent and noisy.
func BenchmarkCall(b *testing.B) {
	f := func() {}
	b.ReportAllocs()
	for b.Loop() {
		mainthread.Call(f)
	}
}
