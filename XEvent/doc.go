// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XEvent 是一个轻量级的事件管理器，支持多重监听、单次回调和批量通知等功能。

功能特性

  - 多重监听：可配置是否允许同一事件注册多个回调
  - 单次执行：可设置回调函数仅执行一次后自动注销

使用手册

1. 事件管理

1.1 创建管理器

	// 创建一个支持多重注册的事件管理器
	mgr := XEvent.NewManager(true)

	// 获取共享的事件管理器
	mgr := XEvent.Shared()

1.2 注册事件

	// 注册事件处理程序
	mgr.Reg(1001, func(args ...any) {
		// 处理事件
	})

	// 注册只执行一次的处理程序
	mgr.Reg(1001, func(args ...any) {
		// 处理事件
	}, true)

1.3 注销事件

	// 注销指定的事件处理程序
	mgr.Unreg(1001, handler)

	// 注销事件的所有处理程序
	mgr.Unreg(1001)

1.4 通知事件

	// 通知事件（不带参数）
	mgr.Notify(1001)

	// 通知事件（带参数）
	mgr.Notify(1001, "param1", 123)

更多信息请参考模块文档。
*/
package XEvent
