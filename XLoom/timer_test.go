// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/eframework-org/GO.UTIL/XTime"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	setup(XPrefs.New().Set(prefsCount, 2).Set(prefsStep, 10).Set(prefsQueue, 1000))

	t.Run("Pool", func(t *testing.T) {
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

	t.Run("Timeout", func(t *testing.T) {
		done := make(chan struct{})
		tm := XTime.GetMillisecond()
		dt := 0
		tm1 := SetTimeout(func() {
			dt = XTime.GetMillisecond() - tm
			close(done)
		}, 500, 0)

		clear := true
		tm2 := SetTimeout(func() { clear = false }, 500, 0)
		ClearTimeout(tm2, 0)

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("定时器回调超时")
		}

		assert.Greater(t, tm1, 0, "返回的定时器 ID 应该为正数")
		assert.GreaterOrEqual(t, dt, 500, "等待时间应当大于等于 500 毫秒")
		assert.Equal(t, true, clear, "清除的定时器不应当被回调")

		assert.Equal(t, -1, SetTimeout(nil, 100, 0), "传入空的回调函数应当返回 -1")
		assert.Equal(t, -1, SetTimeout(func() {}, -1, 0), "传入小于零的超时时长应当返回 -1")
		assert.Equal(t, -1, SetTimeout(func() {}, 100, -1), "传入非法的 loomID 应当返回 -1")
		assert.Equal(t, -1, SetTimeout(func() {}, 100, 999), "传入越界的 loomID 应当返回 -1")
	})

	t.Run("Interval", func(t *testing.T) {
		count := 0
		done := make(chan struct{})

		tm := XTime.GetMillisecond()
		dt := 0
		tm1 := 0
		tm1 = SetInterval(func() {
			count++
			if count >= 3 {
				dt = XTime.GetMillisecond() - tm
				ClearInterval(tm1, 1)
				close(done)
			}
			panic("test interval panic") // 触发 panic，下一个周期的定时器应当继续执行
		}, 200, 1)

		clear := true
		tm2 := SetInterval(func() { clear = false }, 200, 1)
		ClearInterval(tm2, 1)

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("定时器回调超时")
		}

		assert.Greater(t, tm1, 0, "返回的定时器 ID 应该为正数")
		assert.GreaterOrEqual(t, dt, 600, "等待时间应当大于等于 600 毫秒")
		assert.Equal(t, true, clear, "清除的定时器不应当被回调")

		assert.Equal(t, -1, SetInterval(nil, 100, 0), "传入空的回调函数应当返回 -1")
		assert.Equal(t, -1, SetInterval(func() {}, -1, 0), "传入小于零的超时时长应当返回 -1")
		assert.Equal(t, -1, SetInterval(func() {}, 100, -1), "传入非法的 loomID 应当返回 -1")
		assert.Equal(t, -1, SetInterval(func() {}, 100, 999), "传入越界的 loomID 应当返回 -1")
	})
}
