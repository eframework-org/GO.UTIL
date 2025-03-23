// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XApp

import (
	"sync"
)

// IBase 定义了应用程序的生命周期接口。
// 实现此接口的类型可以通过 Run 函数启动并由框架管理其生命周期。
type IBase interface {
	// Awake 在应用程序启动前调用，用于进行初始化检查。
	// 返回 false 将导致应用程序终止启动。
	Awake() bool

	// Start 在应用程序启动时调用，用于执行初始化操作。
	Start()

	// Stop 在应用程序退出时调用，用于执行清理操作。
	// wait 参数用于同步等待清理完成。
	Stop(wait *sync.WaitGroup)
}

// Base 是一个泛型基类，提供了 IBase 接口的默认实现。
// 类型参数 T 通常是实现类的类型。
type Base[T any] struct{ this IBase }

// Ctor 是 Base 的构造函数，用于初始化基类。
// obj 参数必须是实现了 IBase 接口的实例。
func (app *Base[T]) Ctor(obj any) { app.this = obj.(IBase) }

// Awake 的默认实现，返回 true 表示初始化检查通过。
func (app *Base[T]) Awake() bool { return true }

// Start 的默认实现，无操作。
func (app *Base[T]) Start() {}

// Stop 的默认实现，无操作。
func (app *Base[T]) Stop(wait *sync.WaitGroup) {}
