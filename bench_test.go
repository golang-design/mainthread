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
