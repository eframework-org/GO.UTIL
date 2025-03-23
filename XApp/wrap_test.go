// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XApp

import (
	"sync"
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XLog"
	"github.com/stretchr/testify/assert"
)

type MockApp struct {
	started bool
	stopped bool
	mu      sync.Mutex
}

func (app *MockApp) Awake() bool {
	return true
}

func (app *MockApp) Start() {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.started = true
	XLog.Notice("MockApp started.")
}

func (app *MockApp) Stop(wg *sync.WaitGroup) {
	app.mu.Lock()
	defer app.mu.Unlock()
	wg.Add(1)       // 先增加计数
	defer wg.Done() // 确保在函数返回时减少计数
	app.stopped = true
	XLog.Notice("MockApp stopped.")
}

func TestRunAndQuit(t *testing.T) {
	// 替换退出函数以避免 os.Exit
	originalExitFunc := exitFunc
	exitFunc = func(code int) {}
	defer func() { exitFunc = originalExitFunc }()

	mockApp := &MockApp{}
	done := make(chan struct{})

	go func() {
		Run(mockApp)
		close(done)
	}()

	// 等待应用程序启动
	time.Sleep(100 * time.Millisecond)

	mockApp.mu.Lock()
	assert.True(t, mockApp.started, "App should be started")
	mockApp.mu.Unlock()

	// 触发退出
	Quit()

	// 等待应用程序完全停止
	select {
	case <-done:
		// 应用程序已经停止
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for app to stop")
	}

	mockApp.mu.Lock()
	assert.True(t, mockApp.stopped, "App should be stopped")
	mockApp.mu.Unlock()
}
