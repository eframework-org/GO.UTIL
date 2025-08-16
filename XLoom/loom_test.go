// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"fmt"
	"sync"
	"testing"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/eframework-org/GO.UTIL/XTime"
)

func TestLoom(t *testing.T) {
	t.Run("Count", func(t *testing.T) {
		setup(XPrefs.New().Set(prefsCount, 1).Set(prefsStep, 10).Set(prefsQueue, 100000))

		wg := sync.WaitGroup{}
		for range 1000 {
			wg.Add(1)
			RunIn(func() {
				defer wg.Done()

				fmt.Println(fmt.Sprintf("%v-%v", XTime.GetMillisecond(), "111"))
			}, 0)
		}

		wg.Wait()
	})

	// time.Sleep()

	// t.Run("ID", func(t *testing.T) {
	// 	setup(XPrefs.Asset())

	// 	var wg sync.WaitGroup
	// 	wg.Add(1)

	// 	RunIn(func() {
	// 		// 测试当前运行的 loom ID
	// 		lid := ID()
	// 		assert.Equal(t, 0, lid, "Should be running in loom 0")
	// 		wg.Done()
	// 	}, 0)

	// 	wg.Wait()
	// })

	// t.Run("RunIn", func(t *testing.T) {
	// 	setup(XPrefs.Asset())

	// 	var wg sync.WaitGroup
	// 	executed := false
	// 	wg.Add(1)

	// 	RunIn(func() {
	// 		executed = true
	// 		wg.Done()
	// 	}, 0)

	// 	wg.Wait()
	// 	assert.True(t, executed, "Task should be executed")
	// })

	// t.Run("Pause/Resume", func(t *testing.T) {
	// 	setup(XPrefs.Asset())

	// 	var wg sync.WaitGroup

	// 	Pause(0)
	// 	wg.Add(1)
	// 	RunIn(func() {
	// 		wg.Done()
	// 	}, 0)
	// 	Resume(0)

	// 	done := make(chan struct{})
	// 	go func() {
	// 		wg.Wait()
	// 		close(done)
	// 	}()

	// 	select {
	// 	case <-done:
	// 	case <-time.After(time.Second):
	// 		t.Fatal("Task did not execute after resume")
	// 	}
	// })

	// t.Run("Metrics", func(t *testing.T) {
	// 	setup(XPrefs.Asset())

	// 	// 1. 首先测试正常运行时的指标
	// 	taskCount := 100
	// 	for range taskCount {
	// 		RunIn(func() {}, 0)
	// 	}
	// 	expectedFPS := 100 // 因为step=10ms，所以理论上每秒应该有100帧
	// 	expectedQPS := 100

	// 	// 等待一个完整的统计周期
	// 	time.Sleep(time.Millisecond * 1200)

	// 	// 验证正常运行时的指标
	// 	assert.InDelta(t, expectedFPS, FPS(0), float64(expectedFPS)*0.3, "FPS should be around %d (±30%%)", expectedFPS)
	// 	assert.InDelta(t, expectedFPS, testutil.ToFloat64(loomFPSGauges[0]), float64(expectedFPS)*0.3, "FPS should be around %d (±30%%)", expectedFPS)

	// 	assert.InDelta(t, expectedQPS, QPS(0), float64(expectedQPS)*0.3, "QPS should be around %d (±30%%)", expectedQPS)
	// 	assert.InDelta(t, expectedQPS, testutil.ToFloat64(loomQPSGauges[0]), float64(expectedQPS)*0.3, "QPS should be around %d (±30%%)", expectedQPS)

	// 	assert.Equal(t, 100, int(testutil.ToFloat64(loomQueryCounters[0])), "Query count should be 100")
	// 	assert.Equal(t, 100, int(testutil.ToFloat64(loomQueryCounter)), "Total query count should be 100")

	// 	// 2. 测试暂停状态下的指标
	// 	Pause(0)

	// 	// 尝试发送任务
	// 	for range taskCount {
	// 		RunIn(func() {}, 0)
	// 	}

	// 	// 等待一个完整的统计周期
	// 	time.Sleep(time.Millisecond * 1200)

	// 	// 验证暂停时的指标
	// 	assert.Zero(t, FPS(0), "FPS should be 0 while paused")
	// 	assert.Zero(t, testutil.ToFloat64(loomFPSGauges[0]), "FPS should be 0 while paused")

	// 	assert.Zero(t, QPS(0), "QPS should be 0 while paused")
	// 	assert.Zero(t, testutil.ToFloat64(loomQPSGauges[0]), "QPS should be 0 while paused")

	// 	assert.Equal(t, 100, int(testutil.ToFloat64(loomQueryCounters[0])), "Query count should be 100")
	// 	assert.Equal(t, 100, int(testutil.ToFloat64(loomQueryCounter)), "Total query count should be 100")

	// 	// 3. 测试恢复后的指标
	// 	Resume(0)

	// 	// 等待一个完整的统计周期
	// 	time.Sleep(time.Millisecond * 1200)

	// 	// 验证恢复后的指标
	// 	assert.InDelta(t, expectedFPS, FPS(0), float64(expectedFPS)*0.3, "FPS should be around %d (±30%%)", expectedFPS)
	// 	assert.InDelta(t, expectedFPS, testutil.ToFloat64(loomFPSGauges[0]), float64(expectedFPS)*0.3, "FPS should be around %d (±30%%)", expectedFPS)

	// 	assert.InDelta(t, expectedQPS, QPS(0), float64(expectedQPS)*0.3, "QPS should be around %d (±30%%)", expectedQPS)
	// 	assert.InDelta(t, expectedQPS, testutil.ToFloat64(loomQPSGauges[0]), float64(expectedQPS)*0.3, "QPS should be around %d (±30%%)", expectedQPS)

	// 	assert.Equal(t, 200, int(testutil.ToFloat64(loomQueryCounters[0])), "Query count should be 200")
	// 	assert.Equal(t, 200, int(testutil.ToFloat64(loomQueryCounter)), "Total query count should be 200")

	// 	// 4. 测试无效的处理器ID
	// 	assert.Equal(t, 0, FPS(-1), "Invalid PID should return 0")
	// 	assert.Equal(t, 0, FPS(999), "Out of range PID should return 0")
	// 	assert.Equal(t, 0, QPS(-1), "Invalid PID should return 0")
	// 	assert.Equal(t, 0, QPS(999), "Out of range PID should return 0")
	// })
}
