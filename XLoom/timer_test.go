// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	setup(XPrefs.New().Set(prefsCount, 2).Set(prefsStep, 10).Set(prefsQueue, 1000))

	t.Run("Timer Pool", func(t *testing.T) {
		// 测试定时器对象池
		timer1 := timerPool.Get().(*timer)
		assert.NotNil(t, timer1)
		timer1.reset()
		timerPool.Put(timer1)

		timer2 := timerPool.Get().(*timer)
		assert.NotNil(t, timer2)
		assert.Equal(t, 0, timer2.id)
		assert.Nil(t, timer2.callback)
	})

	t.Run("SetTimeout Basic", func(t *testing.T) {
		var executed int32
		done := make(chan struct{})

		// 在处理器 0 上设置一个超时调用
		id := SetTimeout(func() {
			atomic.AddInt32(&executed, 1)
			close(done)
		}, 50, 0)

		assert.Greater(t, id, 0, "Timer ID should be positive")

		select {
		case <-done:
			assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
		case <-time.After(time.Second):
			t.Fatal("Timeout did not execute in time")
		}
	})

	t.Run("SetTimeout Invalid", func(t *testing.T) {
		// 测试无效参数
		assert.Equal(t, -1, SetTimeout(nil, 100, 0), "Nil callback should return -1")
		assert.Equal(t, -1, SetTimeout(func() {}, -1, 0), "Negative timeout should return -1")
		assert.Equal(t, -1, SetTimeout(func() {}, 100, -1), "Invalid PID should return -1")
		assert.Equal(t, -1, SetTimeout(func() {}, 100, 999), "Out of range PID should return -1")
	})

	t.Run("SetInterval Basic", func(t *testing.T) {
		var count int32
		done := make(chan struct{})

		// 设置一个间歇调用
		id := SetInterval(func() {
			if atomic.AddInt32(&count, 1) >= 3 {
				close(done)
			}
		}, 50, 0)

		assert.Greater(t, id, 0, "Timer ID should be positive")

		select {
		case <-done:
			assert.GreaterOrEqual(t, atomic.LoadInt32(&count), int32(3))
		case <-time.After(time.Second):
			t.Fatal("Interval did not execute enough times")
		}

		// 清除间歇调用
		ClearInterval(id, 0)
	})

	t.Run("SetInterval Invalid", func(t *testing.T) {
		// 测试无效参数
		assert.Equal(t, -1, SetInterval(nil, 100, 0), "Nil callback should return -1")
		assert.Equal(t, -1, SetInterval(func() {}, -1, 0), "Negative interval should return -1")
		assert.Equal(t, -1, SetInterval(func() {}, 100, -1), "Invalid PID should return -1")
		assert.Equal(t, -1, SetInterval(func() {}, 100, 999), "Out of range PID should return -1")
	})

	t.Run("Timer Panic Recovery", func(t *testing.T) {
		var executed int32

		// 测试定时器中的 panic 恢复
		SetTimeout(func() {
			atomic.AddInt32(&executed, 1)
			panic("test panic")
		}, 50, 0)

		// 等待足够时间让定时器执行
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	})

	t.Run("Timer Clear", func(t *testing.T) {
		var executed int32

		// 设置一个定时器并立即清除
		id := SetTimeout(func() {
			atomic.AddInt32(&executed, 1)
		}, 50, 0)

		ClearTimeout(id, 0)

		// 等待一段时间确保定时器被清除
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, int32(0), atomic.LoadInt32(&executed))
	})

	t.Run("Timer Update", func(t *testing.T) {
		// 测试定时器更新逻辑
		var wg sync.WaitGroup
		wg.Add(1)

		SetTimeout(func() {
			defer wg.Done()
		}, 50, 0)

		// 手动触发更新
		updateTimer(0, 60)

		wg.Wait()
	})
}
