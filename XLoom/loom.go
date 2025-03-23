// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eframework-org/GO.UTIL/XLog"
	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/eframework-org/GO.UTIL/XTime"
	"github.com/illumitacit/gostd/quit"
	"github.com/petermattis/goid"
)

const (
	prefsCount        = "Loom/Count" // 线程数量配置键，用于设置线程池大小
	prefsCountDefault = 1            // 默认线程数量，当未配置时使用此值
	prefsStep         = "Loom/Step"  // 更新步长配置键，用于控制线程更新频率（毫秒）
	prefsStepDefault  = 10           // 默认更新步长，当未配置时使用此值
	prefsQueue        = "Loom/Queue" // 队列大小配置键，用于设置每个线程的任务队列容量
	prefsQueueDefault = 50000        // 默认队列大小，当未配置时使用此值
)

var (
	initMu     sync.Mutex            // 初始化互斥锁，用于保护初始化过程
	initSigMap sync.Map              //map[int64]chan os.Signal
	closeWait  sync.WaitGroup        // 等待所有提交处理器完成
	loomIDMap  = make(map[int64]int) // 线程映射表，用于存储 goroutine ID 到 loom ID 的映射关系
	loomIDMu   sync.Mutex            // 线程映射表互斥锁，用于保护映射表的并发访问
)

var (
	loomCount int           // 线程总数，表示当前运行的线程数量
	loomTask  []chan func() // 线程任务队列，每个线程一个独立的任务通道
	loomPause []bool        // 线程暂停状态，true 表示暂停，false 表示运行
	loomFps   []int         // 线程刷新帧率统计，记录每个线程的每秒帧数
	loomQps   []int         // 线程处理速率统计，记录每个线程的每秒处理任务数
)

func init() { setup(XPrefs.Asset()) }

// setup 初始化线程系统。
func setup(prefs XPrefs.IBase) {
	if prefs == nil {
		XLog.Panic("XLoom.Init: prefs is nil.")
		return
	}

	initMu.Lock()
	defer initMu.Unlock()

	count := prefs.GetInt(prefsCount, prefsCountDefault)
	step := prefs.GetInt(prefsStep, prefsStepDefault)
	queue := prefs.GetInt(prefsQueue, prefsQueueDefault)

	if count <= 0 || step <= 0 || queue <= 0 {
		XLog.Panic("XLoom.Init: invalid parameters - count: %v, step: %v, queue: %v",
			count, step, queue)
		return
	}

	loomCount = count
	loomTask = make([]chan func(), count)
	loomPause = make([]bool, count)
	loomFps = make([]int, count)
	loomQps = make([]int, count)

	for i := range count {
		loomTask[i] = make(chan func(), queue)
	}

	setupTimer(count)

	wg := sync.WaitGroup{}
	for i := range count {
		wg.Add(1)
		doneOnce := sync.Once{}
		RunAsyncT1(func(pid int) {
			initVal, _ := initSigMap.LoadOrStore(pid, make(chan os.Signal, 1))
			ch := initVal.(chan os.Signal)
			signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
			closeWait.Add(1)
			quit.GetWaiter().Add(1)
			defer func() {
				quit.GetWaiter().Done()
				closeWait.Done()
			}()

			loomIDMu.Lock()
			loomIDMap[goid.Get()] = pid
			loomIDMu.Unlock()

			updateTicker := time.NewTicker(time.Millisecond * time.Duration(step))
			defer updateTicker.Stop()

			doneOnce.Do(func() { // 确保只调用一次，否则recover后会重复调用
				wg.Done() // 确保线程启动完成
			})

			lastTime := XTime.GetMillisecond()
			frameCount := 0
			queryCount := 0

			for {
				if loomPause[pid] {
					select {
					case <-updateTicker.C:
						// 在暂停状态下重置计数器和指标
						frameCount = 0
						queryCount = 0
						loomFps[pid] = 0
						loomQps[pid] = 0
						lastTime = XTime.GetMillisecond() // 更新时间戳，避免恢复后的突然跳变
					case sig, ok := <-ch:
						if ok {
							XLog.Notice("XLoom.Loop(%v): receive signal of %v.", pid, sig.String())
						} else {
							XLog.Notice("XLoom.Loop(%v): channel of signal is closed.", pid)
						}
						return
					case <-quit.GetQuitChannel():
						XLog.Notice("XLoom.Loop(%v): receive signal of QUIT.", pid)
						return
					}
				} else {
					nowTime := XTime.GetMillisecond()
					deltaTime := nowTime - lastTime

					if deltaTime >= 1000 {
						fps := int(float64(frameCount) * 1000 / float64(deltaTime))
						qps := int(float64(queryCount) * 1000 / float64(deltaTime))
						loomFps[pid] = fps
						loomQps[pid] = qps
						frameCount = 0
						queryCount = 0
						lastTime = nowTime
					}

					select {
					case runIn, ok := <-loomTask[pid]:
						if ok {
							queryCount++
							runIn()
						} else {
							XLog.Error("XLoom.Loop(%v): get runin with ret false.", pid)
						}
					case <-updateTicker.C:
						frameCount++
						updateTimer(pid, deltaTime)
					case sig, ok := <-ch:
						if ok {
							XLog.Notice("XLoom.Loop(%v): receive signal of %v.", pid, sig.String())
						} else {
							XLog.Notice("XLoom.Loop(%v): channel of signal is closed.", pid)
						}
						return
					case <-quit.GetQuitChannel():
						XLog.Notice("XLoom.Loop(%v): receive signal of QUIT.", pid)
						return
					}
				}
			}
		}, i, true)
	}

	XLog.Notice("XLoom.Init: allocated %v loom(s).", count)
	loomCount = count
	wg.Wait()
}

// Pause 暂停指定线程或所有线程。
// loomID 为可选的目标线程 ID，如果未指定，则暂停所有线程。
func Pause(loomID ...int) {
	if len(loomID) == 1 {
		lid := loomID[0]
		if lid < 0 {
			XLog.Critical("XLoom.Pause: loom id of %v can not be zero or negative.", lid)
			return
		}
		if lid >= loomCount {
			XLog.Critical("XLoom.Pause: loom id of %v can not equals or greater than: %v", lid, Count())
			return
		}
		loomPause[lid] = true
	} else {
		for i := range loomPause {
			loomPause[i] = true
		}
	}
}

// Resume 恢复指定线程或所有线程。
// loomID 为可选的目标线程 ID，如果未指定，则恢复所有线程。
func Resume(loomID ...int) {
	if len(loomID) == 1 {
		lid := loomID[0]
		if lid < 0 {
			XLog.Critical("XLoom.Resume: loom id of %v can not be zero or negative.", lid)
			return
		}
		if lid >= loomCount {
			XLog.Critical("XLoom.Resume: loom id of %v can not equals or greater than: %v", lid, Count())
			return
		}
		loomPause[lid] = false
	} else {
		for i := range loomPause {
			loomPause[i] = false
		}
	}
}

// RunIn 在指定线程中执行任务。
// callback 为要执行的任务函数。
// loomID 为可选的目标线程 ID，如果未指定，默认在线程 0 中执行。
func RunIn(callback func(), loomID ...int) {
	if callback == nil {
		XLog.Critical("XLoom.RunIn: callback can not be nil.")
		return
	}
	lid := -1
	if len(loomID) == 1 {
		lid = loomID[0]
	} else {
		lid = 0
	}
	if lid < 0 {
		XLog.Critical("XLoom.RunIn: loom id of %v can not be zero or negative.", lid)
		return
	}
	if lid >= loomCount {
		XLog.Critical("XLoom.RunIn: loom id of %v can not equals or greater than: %v", lid, Count())
		return
	}
	ch := loomTask[lid]
	select {
	case ch <- callback:
	default:
		XLog.Critical("XLoom.RunIn: too many runins of %v", lid)
	}
}

// Count 返回线程总数。
func Count() int { return loomCount }

// ID 获取当前 goroutine 所在的 loom ID。
// 如果指定了 goroutineID，则返回该线程的线程 ID。
// 如果线程未绑定线程，返回 -1。
func ID(goroutineID ...int64) int {
	// TONOTICE: 不使用sync.Map避免引起值类型的装箱和拆箱
	// 尽量在业务线程中调用，业务线程之外调用可能存在并发读写问题
	var tgid int64
	if len(goroutineID) == 1 {
		tgid = goroutineID[0]
	} else {
		tgid = goid.Get()
	}
	if pid, ok := loomIDMap[tgid]; ok {
		return pid
	}
	return -1
}

// FPS 获取指定线程的刷新帧率。
// loomID 为可选的目标线程 ID，如果未指定，返回当前线程的刷新帧率。
// 返回每秒帧数，如果线程 ID 无效则返回 0。
func FPS(loomID ...int) int {
	lid := -1
	if len(loomID) == 1 {
		lid = loomID[0]
	} else {
		lid = ID()
	}
	if lid < 0 {
		XLog.Critical("XLoom.FPS: loom id of %v can not be zero or negative.", lid)
		return 0
	}
	if lid >= loomCount {
		XLog.Critical("XLoom.FPS: loom id of %v can not equals or greater than: %v", lid, Count())
		return 0
	}
	return loomFps[lid]
}

// QPS 获取指定线程的处理速率。
// loomID 为可选的目标线程 ID，如果未指定，返回当前线程的处理速率。
// 返回每秒处理的任务数，如果线程 ID 无效则返回 0。
func QPS(loomID ...int) int {
	lid := -1
	if len(loomID) == 1 {
		lid = loomID[0]
	} else {
		lid = ID()
	}
	if lid < 0 {
		XLog.Critical("XLoom.QPS: loom id of %v can not be zero or negative.", lid)
		return 0
	}
	if lid >= loomCount {
		XLog.Critical("XLoom.QPS: loom id of %v can not equals or greater than: %v", lid, Count())
		return 0
	}
	return loomQps[lid]
}
