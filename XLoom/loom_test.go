// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/stretchr/testify/assert"
)

func TestLoom(t *testing.T) {
	//注：避免init()函数重复setup
	initSigMap.Range(func(key any, value any) bool {
		if ch, ok := value.(chan os.Signal); ok {
			signal.Stop(ch)
			close(ch)
		}
		return true
	})
	closeWait.Wait()
	initSigMap = sync.Map{}

	setup(XPrefs.New().Set(prefsCount, 2).Set(prefsStep, 10).Set(prefsQueue, 1000))

	t.Run("Loom Count", func(t *testing.T) {
		// 测试当前活动的 loom 数量
		assert.Equal(t, 2, Count(), "Should have 2 looms")
	})

	t.Run("Loom ID", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		RunIn(func() {
			// 测试当前运行的 loom ID
			lid := ID()
			assert.Equal(t, 0, lid, "Should be running in loom 0")
			wg.Done()
		}, 0)

		wg.Wait()
	})

	t.Run("Task Execution", func(t *testing.T) {
		var wg sync.WaitGroup
		executed := false
		wg.Add(1)

		RunIn(func() {
			executed = true
			wg.Done()
		}, 0)

		wg.Wait()
		assert.True(t, executed, "Task should be executed")
	})

	t.Run("Pause Resume", func(t *testing.T) {
		var wg sync.WaitGroup

		Pause(0)
		wg.Add(1)
		RunIn(func() {
			wg.Done()
		}, 0)
		Resume(0)

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Task did not execute after resume")
		}
	})

	t.Run("Metrics", func(t *testing.T) {
		// 1. 首先测试正常运行时的指标
		taskCount := 100
		for range taskCount {
			RunIn(func() {}, 0)
		}

		// 等待一个完整的统计周期
		time.Sleep(1200 * time.Millisecond)

		// 验证正常运行时的指标
		fps := FPS(0)
		expectedFPS := 100 // 因为step=10ms，所以理论上每秒应该有100帧
		assert.InDelta(t, expectedFPS, fps, float64(expectedFPS)*0.3, "FPS should be around %d (±30%%)", expectedFPS)

		qps := QPS(0)
		assert.InDelta(t, taskCount, qps, float64(taskCount)*0.3, "QPS should be around %d (±30%%)", taskCount)

		// 2. 测试暂停状态下的指标
		Pause(0)

		// 尝试发送任务
		for i := 0; i < taskCount; i++ {
			RunIn(func() {}, 0)
		}

		// 等待一个完整的统计周期
		time.Sleep(1200 * time.Millisecond)

		// 验证暂停时的指标
		pauseFPS := FPS(0)
		pauseQPS := QPS(0)
		assert.Zero(t, pauseFPS, "FPS should be 0 while paused")
		assert.Zero(t, pauseQPS, "QPS should be 0 while paused")

		// 3. 测试恢复后的指标
		Resume(0)

		// 等待一个完整的统计周期
		time.Sleep(1200 * time.Millisecond)

		// 验证恢复后的指标
		fps = FPS(0)
		expectedFPS = 100 // 因为step=10ms，所以理论上每秒应该有100帧
		assert.InDelta(t, expectedFPS, fps, float64(expectedFPS)*0.3, "FPS should be around %d (±30%%)", expectedFPS)

		qps = QPS(0)
		assert.InDelta(t, taskCount, qps, float64(taskCount)*0.3, "QPS should be around %d (±30%%)", taskCount)

		// 4. 测试无效的处理器ID
		assert.Equal(t, 0, FPS(-1), "Invalid PID should return 0")
		assert.Equal(t, 0, FPS(999), "Out of range PID should return 0")
		assert.Equal(t, 0, QPS(-1), "Invalid PID should return 0")
		assert.Equal(t, 0, QPS(999), "Out of range PID should return 0")
	})
}
