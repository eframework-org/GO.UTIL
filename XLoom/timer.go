// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"sync"
	"sync/atomic"

	"github.com/eframework-org/GO.UTIL/XLog"
)

var (
	timerPool = sync.Pool{New: func() any {
		obj := new(timer)
		return obj
	}}
	timerIID    int64      // 定时器自增标识
	allTimers   [][]*timer // 所有定时器
	newTimers   [][]*timer // 新的定时器
	newTimersLk []sync.Mutex
	delTimers   [][]int // 待删除的定时器
	delTimersLk []sync.Mutex
)

// timer 定义了一个定时器的基本结构。
type timer struct {
	id       int    // 定时器唯一标识
	callback func() // 定时器触发时执行的回调函数
	period   int    // 定时器周期（毫秒），用于重复执行的间歇时间
	tick     int    // 当前剩余时间（毫秒），倒计时到 0 时触发回调
	repeat   bool   // 是否重复执行，true 表示间歇调用，false 表示超时调用
	panic    bool   // 是否发生异常，用于异常恢复控制
}

// reset 重置定时器到初始状态。
func (tm *timer) reset() *timer {
	tm.id = 0
	tm.callback = nil
	tm.period = 0
	tm.tick = 0
	tm.repeat = false
	tm.panic = false
	return tm
}

// setupTimer 初始化定时器系统。
func setupTimer(num int) {
	allTimers = make([][]*timer, num)
	newTimers = make([][]*timer, num)
	newTimersLk = make([]sync.Mutex, num)
	delTimers = make([][]int, num)
	delTimersLk = make([]sync.Mutex, num)
}

// updateTimer 更新指定线成的定时器状态。
func updateTimer(pid int, delta int) {
	if newTimers[pid] != nil && len(newTimers[pid]) > 0 {
		newTimersLk[pid].Lock()
		allTimers[pid] = append(allTimers[pid], newTimers[pid]...)
		newTimers[pid] = newTimers[pid][:0]
		newTimersLk[pid].Unlock()
	}
	if delTimers[pid] != nil && len(delTimers[pid]) > 0 {
		delTimersLk[pid].Lock()
		for _, id := range delTimers[pid] {
			for idx, timer := range allTimers[pid] {
				if id == timer.id {
					allTimers[pid] = append(allTimers[pid][:idx], allTimers[pid][idx+1:]...)
					timerPool.Put(timer.reset())
					break
				}
			}
		}
		delTimers[pid] = delTimers[pid][:0]
		delTimersLk[pid].Unlock()
	}
	if allTimers[pid] != nil {
		for _, timer := range allTimers[pid] {
			timer.tick -= delta
			if timer.panic {
				if timer.repeat { // interval 发生 panic 不取消定时器
					timer.panic = false
					timer.tick = timer.period
				} else { // timeout 发生 panic 则直接移除
					delTimersLk[pid].Lock()
					delTimers[pid] = append(delTimers[pid], timer.id)
					delTimersLk[pid].Unlock()
					continue
				}
			}
			if timer.tick <= 0 { // 因存在固定刷新间歇，可能会导致间歇调用的周期越来越长
				if timer.callback != nil {
					timer.panic = true
					timer.callback()
					timer.panic = false
				}
				if !timer.repeat {
					delTimersLk[pid].Lock()
					delTimers[pid] = append(delTimers[pid], timer.id)
					delTimersLk[pid].Unlock()
				} else {
					timer.tick = timer.period
				}
			}
		}
	}
}

// SetTimeout 设置一个超时调用。
// callback 为要执行的回调函数。
// timeout 为超时时间（毫秒）。
// loomID 为可选的目标线程 ID，如果未指定，在当前线程中执行。
// 返回定时器 ID，如果参数无效则返回 -1。
func SetTimeout(callback func(), timeout int, loomID ...int) int {
	if callback == nil {
		XLog.Critical("XLoom.SetTimeout: callback can not be nil.")
		return -1
	}
	if timeout < 0 {
		XLog.Critical("XLoom.SetTimeout: timeout of %v can not be zero or negative.", timeout)
		return -1
	}
	lid := -1
	if len(loomID) == 1 {
		lid = loomID[0]
	} else {
		lid = ID()
	}
	if lid < 0 {
		XLog.Critical("XLoom.SetTimeout: loom id of %v can not be zero or negative.", lid)
		return -1
	}
	if lid >= loomCount {
		XLog.Critical("XLoom.SetTimeout: loom id of %v can not equals or greater than: %v", lid, Count())
		return -1
	}

	timer := timerPool.Get().(*timer)
	timer.id = int(atomic.AddInt64(&timerIID, 1))
	timer.callback = callback
	timer.tick = timeout
	timer.period = timer.tick
	timer.repeat = false

	newTimersLk[lid].Lock()
	newTimers[lid] = append(newTimers[lid], timer)
	newTimersLk[lid].Unlock()
	return timer.id
}

// ClearTimeout 取消一个超时调用。
// id 为要取消的定时器 ID。
// loomID 为可选的目标线程 ID，如果未指定，在当前线程中执行。
func ClearTimeout(id int, loomID ...int) {
	lid := -1
	if len(loomID) == 1 {
		lid = loomID[0]
	} else {
		lid = ID()
	}
	if lid < 0 {
		XLog.Critical("XLoom.ClearTimeout: loom id of %v can not be zero or negative.", lid)
		return
	}
	if lid >= loomCount {
		XLog.Critical("XLoom.ClearTimeout: loom id of %v can not equals or greater than: %v", lid, Count())
		return
	}

	delTimersLk[lid].Lock()
	delTimers[lid] = append(delTimers[lid], id)
	delTimersLk[lid].Unlock()
}

// SetInterval 设置一个间歇调用。
// callback 为要执行的回调函数。
// interval 为调用间歇（毫秒）。
// loomID 为可选的目标线程 ID，如果未指定，在当前线程中执行。
// 返回定时器 ID，如果参数无效则返回 -1。
func SetInterval(callback func(), interval int, loomID ...int) int {
	if callback == nil {
		XLog.Critical("XLoom.SetInterval: callback can not be nil.")
		return -1
	}
	if interval < 0 {
		XLog.Critical("XLoom.SetInterval: interval of %v can not be zero or negative.", interval)
		return -1
	}
	lid := -1
	if len(loomID) == 1 {
		lid = loomID[0]
	} else {
		lid = ID()
	}
	if lid < 0 {
		XLog.Critical("XLoom.SetInterval: loom id of %v can not be zero or negative.", lid)
		return -1
	}
	if lid >= loomCount {
		XLog.Critical("XLoom.SetInterval: loom id of %v can not equals or greater than: %v", lid, Count())
		return -1
	}

	timer := timerPool.Get().(*timer)
	timer.id = int(atomic.AddInt64(&timerIID, 1))
	timer.callback = callback
	timer.tick = interval
	timer.period = timer.tick
	timer.repeat = true

	newTimersLk[lid].Lock()
	newTimers[lid] = append(newTimers[lid], timer)
	newTimersLk[lid].Unlock()
	return timer.id
}

// ClearInterval 取消一个间歇调用。
// id 为要取消的定时器 ID。
// loomID 为可选的目标线程 ID，如果未指定，在当前线程中执行。
func ClearInterval(id int, loomID ...int) { ClearTimeout(id, loomID...) }
