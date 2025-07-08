// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XApp

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eframework-org/GO.UTIL/XLog"
	"github.com/illumitacit/gostd/quit"
)

var (
	shared   IBase
	runOnce  sync.Once
	quitOnce sync.Once
)

// Shared 返回应用程序的单例实例。
// 类型参数 T 必须是 IBase 接口的实现类型。
// 在应用程序启动前调用将返回 nil。
func Shared[T IBase]() T { return shared.(T) }

// Run 启动并运行应用程序。
// app 参数必须是 IBase 接口的实现实例。
// 此函数会阻塞直到应用程序退出。
// 应用程序可以通过以下方式退出：
//   - 调用 Quit 函数
//   - 接收到 SIGTERM 或 SIGINT 信号
//   - Awake 返回 false
func Run(app IBase) {
	runOnce.Do(func() {
		if app == nil {
			XLog.Panic("XApp.Run: app is nil.")
		}
		shared = app

		if !app.Awake() {
			XLog.Panic("XApp.Run: app awake failed.")
		}
		XLog.Notice("XApp.Run: app has been awaked.")

		app.Start()
		XLog.Notice("XApp.Run: app has been started.")

		defer func() {
			wg := &sync.WaitGroup{}
			app.Stop(wg)
			wg.Wait()
			XLog.Notice("XApp.Run: app has been stopped.")
		}()

		for {
			defer func() {
				quit.GetWaiter().Wait()
			}()
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
			for {
				select {
				case sig, ok := <-ch:
					if ok {
						XLog.Notice("XApp.Listen: receive signal of %v.", sig.String())
					} else {
						XLog.Notice("XApp.Listen: channel of signal is closed.")
					}
					return
				case <-quit.GetQuitChannel():
					XLog.Notice("XApp.Listen: receive signal of quit.")
					return
				}
			}
		}
	})
}

// Quit 触发应用程序退出。
// 此函数可以在任意 goroutine 中安全调用。
// 多次调用只有第一次会生效。
func Quit() { quitOnce.Do(func() { quit.BroadcastShutdown() }) }
