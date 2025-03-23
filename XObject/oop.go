// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XObject

// IThis 定义了对象实例的自引用接口。
// T 为实例的具体类型。
type IThis[T any] interface {
	// This 返回对象实例的指针。
	This() *T
}

// ICtor 定义了对象的无参构造器接口。
type ICtor interface {
	// Ctor 执行对象的初始化操作。
	// obj 为对象实例的指针。
	Ctor(obj any)
}

// ICtorT1 定义了带一个参数的构造器接口。
// T1 为第一个参数的类型。
type ICtorT1[T1 any] interface {
	// CtorT1 执行对象的初始化操作。
	// obj 为对象实例的指针。
	// arg1 为第一个参数。
	CtorT1(obj any, arg1 T1)
}

// ICtorT2 定义了带两个参数的构造器接口。
// T1 为第一个参数的类型。
// T2 为第二个参数的类型。
type ICtorT2[T1 any, T2 any] interface {
	// CtorT2 执行对象的初始化操作。
	// obj 为对象实例的指针。
	// arg1 为第一个参数。
	// arg2 为第二个参数。
	CtorT2(obj any, arg1 T1, arg2 T2)
}

// ICtorT3 定义了带三个参数的构造器接口。
// T1 为第一个参数的类型。
// T2 为第二个参数的类型。
// T3 为第三个参数的类型。
type ICtorT3[T1 any, T2 any, T3 any] interface {
	// CtorT3 执行对象的初始化操作。
	// obj 为对象实例的指针。
	// arg1 为第一个参数。
	// arg2 为第二个参数。
	// arg3 为第三个参数。
	CtorT3(obj any, arg1 T1, arg2 T2, arg3 T3)
}

// New 创建指定类型的对象实例。
// 如果对象实现了 ICtor 接口，则调用其构造函数。
func New[T any]() *T {
	obj := new(T)
	if ictor, ok := any(obj).(ICtor); ok {
		ictor.Ctor(obj)
	}
	return obj
}

// NewT1 创建指定类型的对象实例，并使用一个参数初始化。
// 如果对象实现了 ICtorT1 接口，则调用其构造函数。
func NewT1[T any, T1 any](arg1 T1) *T {
	obj := new(T)
	if ictor, ok := any(obj).(ICtorT1[T1]); ok {
		ictor.CtorT1(obj, arg1)
	}
	return obj
}

// NewT2 创建指定类型的对象实例，并使用两个参数初始化。
// 如果对象实现了 ICtorT2 接口，则调用其构造函数。
func NewT2[T any, T1 any, T2 any](arg1 T1, arg2 T2) *T {
	obj := new(T)
	if ictor, ok := any(obj).(ICtorT2[T1, T2]); ok {
		ictor.CtorT2(obj, arg1, arg2)
	}
	return obj
}

// NewT3 创建指定类型的对象实例，并使用三个参数初始化。
// 如果对象实现了 ICtorT3 接口，则调用其构造函数。
func NewT3[T any, T1 any, T2 any, T3 any](arg1 T1, arg2 T2, arg3 T3) *T {
	obj := new(T)
	if ictor, ok := any(obj).(ICtorT3[T1, T2, T3]); ok {
		ictor.CtorT3(obj, arg1, arg2, arg3)
	}
	return obj
}
