// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XLoom 提供了一个轻量级的任务调度系统，用于管理异步任务、定时器和多线程并发。

功能特性

  - 异步任务：支持执行和异常恢复异步任务
  - 线程管理：支持任务管理、线程暂停/恢复控制、性能监控（FPS/QPS）
  - 定时器管理：支持设置/取消超时和间隔调用

使用手册

1. 异步任务

1.1 基础异步执行

	// 无参数异步执行
	XLoom.RunAsync(func() {
		// 异步代码
	})

	// 带参数异步执行
	XLoom.RunAsyncT1(func(id int) {
		fmt.Printf("Task %d\n", id)
	}, 1)

	// 带异常恢复的异步执行
	XLoom.RunAsync(func() {
		panic("recoverable")
	}, true) // true 表示发生异常时重试

2. 线程管理

2.1 任务调度

	// 在指定线程执行任务
	XLoom.RunIn(func() {
		fmt.Println("在线程0中执行")
	}, 0)

	// 获取当前线程 loom ID
	pid := XLoom.ID()

	// 获取线程性能指标
	fps := XLoom.FPS(0) // 线程0的帧率
	qps := XLoom.QPS(0) // 线程0的处理速率

2.2 线程控制

	// 暂停/恢复单个线程
	XLoom.Pause(0)  // 暂停线程0
	XLoom.Resume(0) // 恢复线程0

	// 暂停/恢复所有线程
	XLoom.Pause()
	XLoom.Resume()

3. 定时器

3.1 超时调用

	// 延迟执行
	id := XLoom.SetTimeout(func() {
		fmt.Println("1秒后执行")
	}, 1000)

	// 取消超时调用
	XLoom.ClearTimeout(id)

3.2 间隔调用

	// 周期执行
	id := XLoom.SetInterval(func() {
		fmt.Println("每秒执行一次")
	}, 1000)

	// 取消间隔调用
	XLoom.ClearInterval(id)

更多信息请参考模块文档。
*/
package XLoom
