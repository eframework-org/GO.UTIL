// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试异步运行功能
func TestRunAsync(t *testing.T) {
	t.Run("Basic Async", func(t *testing.T) {
		done := make(chan struct{})
		RunAsync(func() { close(done) })
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Async function did not complete in time")
		}
	})

	t.Run("Nil Callback", func(t *testing.T) { RunAsync(nil) })

	t.Run("Panic Recovery", func(t *testing.T) {
		done := make(chan struct{})
		var count int32
		RunAsync(func() {
			if atomic.AddInt32(&count, 1) == 1 {
				panic("test panic")
			}
			close(done)
		}, true)
		select {
		case <-done:
			assert.Equal(t, int32(2), atomic.LoadInt32(&count),
				"Function should be executed twice due to recovery")
		case <-time.After(time.Second):
			t.Fatal("Recovery did not complete in time")
		}
	})
}

// 测试带一个参数的异步运行功能
func TestRunAsyncT1(t *testing.T) {
	t.Run("Single Argument", func(t *testing.T) {
		result := make(chan int, 1)
		RunAsyncT1(func(x int) { result <- x * 2 }, 21)
		select {
		case val := <-result:
			assert.Equal(t, 42, val)
		case <-time.After(time.Second):
			t.Fatal("Async function did not complete in time")
		}
	})

	t.Run("Nil Callback", func(t *testing.T) {
		RunAsyncT1(nil, 42)
	})

	t.Run("Panic Recovery", func(t *testing.T) {
		count := 0
		var mu sync.Mutex
		RunAsyncT1(func(x int) {
			mu.Lock()
			count++
			mu.Unlock()
			if count == 1 {
				panic("test panic")
			}
		}, 42, true)
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		assert.Equal(t, 2, count, "Function should be executed twice due to recovery")
		mu.Unlock()
	})
}

// 测试带两个参数的异步运行功能
func TestRunAsyncT2(t *testing.T) {
	t.Run("Two Arguments", func(t *testing.T) {
		result := make(chan string, 1)
		RunAsyncT2(func(x int, y string) { result <- y + ":" + string(rune(x)) }, 65, "A")
		select {
		case val := <-result:
			assert.Equal(t, "A:A", val)
		case <-time.After(time.Second):
			t.Fatal("Async function did not complete in time")
		}
	})

	t.Run("Nil Callback", func(t *testing.T) { RunAsyncT2(nil, 42, "test") })

	t.Run("Panic Recovery", func(t *testing.T) {
		count := 0
		var mu sync.Mutex
		RunAsyncT2(func(x int, y string) {
			mu.Lock()
			count++
			mu.Unlock()
			if count == 1 {
				panic("test panic")
			}
		}, 42, "test", true)
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		assert.Equal(t, 2, count, "Function should be executed twice due to recovery")
		mu.Unlock()
	})
}

// 测试带三个参数的异步运行功能
func TestRunAsyncT3(t *testing.T) {
	t.Run("Three Arguments", func(t *testing.T) {
		result := make(chan string, 1)
		RunAsyncT3(func(x int, y string, z bool) {
			if z {
				result <- y + ":" + string(rune(x))
			}
		}, 65, "A", true)
		select {
		case val := <-result:
			assert.Equal(t, "A:A", val)
		case <-time.After(time.Second):
			t.Fatal("Async function did not complete in time")
		}
	})

	t.Run("Nil Callback", func(t *testing.T) { RunAsyncT3(nil, 42, "test", true) })

	t.Run("Panic Recovery", func(t *testing.T) {
		count := 0
		var mu sync.Mutex
		RunAsyncT3(func(x int, y string, z bool) {
			mu.Lock()
			count++
			mu.Unlock()
			if count == 1 {
				panic("test panic")
			}
		}, 42, "test", true, true)
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		assert.Equal(t, 2, count, "Function should be executed twice due to recovery")
		mu.Unlock()
	})
}
