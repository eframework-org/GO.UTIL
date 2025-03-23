// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XApp 提供了应用程序生命周期管理功能，用于控制应用程序的启动、运行和退出过程。

功能特性

  - 应用程序生命周期管理：提供 Awake、Start、Stop 等生命周期函数
  - 优雅的启动和退出：支持启动前环境检查和退出时的资源清理
  - 泛型单例模式：内置 Base[T] 泛型基类，支持在应用的任意位置获取应用实例

使用手册

1. 生命周期管理

1.1 应用启动

	type MyApp struct {
		XApp.Base[MyApp]
	}
	app := XObject.New[MyApp]()
	XApp.Run(app)

1.2 应用退出

	XApp.Quit()

2. 单例访问

2.1 获取应用实例

	instance := XApp.Shared[*MyApp]()

更多信息请参考模块文档。
*/
package XApp
