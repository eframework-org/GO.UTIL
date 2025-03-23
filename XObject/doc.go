// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XObject 提供了 Go 语言面向对象编程的增强支持，包括对象构造、实例管理和序列化功能。

功能特性

  - 对象构造器：支持泛型的对象构造和初始化
  - 实例指针：支持对象实例的自引用
  - 参数化构造：支持最多三个参数的构造函数
  - JSON 序列化：支持对象与 JSON 的相互转换

使用手册

1. 对象构造

1.1 基础构造

	// 定义结构体
	type MyStruct struct {
		XObject.ICtor // 实现无参构造
		Name string
	}

	// 实现构造函数
	func (ms *MyStruct) Ctor(obj any) {
		ms.Name = "默认名称"
	}

	// 创建实例
	obj := XObject.New[MyStruct]()

1.2 参数化构造

	// 定义结构体
	type MyStruct struct {
		XObject.ICtorT1[string] // 实现带参构造
		Name string
	}

	// 实现构造函数
	func (ms *MyStruct) CtorT1(obj any, name string) {
		ms.Name = name
	}

	// 创建实例
	obj := XObject.NewT1[MyStruct, string]("张三")

1.3 实例自引用

	// 定义结构体
	type MyStruct struct {
		XObject.IThis[MyStruct] // 实现自引用
		this *MyStruct
	}

	// 实现自引用方法
	func (ms *MyStruct) This() *MyStruct {
		return ms.this
	}

	// 在构造函数中初始化
	func (ms *MyStruct) Ctor(obj any) {
		ms.this = obj.(*MyStruct)
	}

2. 序列化

2.1 JSON 转换

	// 对象转 JSON
	json, err := XObject.ToJson(obj)
	json, err := XObject.ToJson(obj, true) // 格式化输出

	// JSON 转对象
	err := XObject.FromJson(json, &obj)

2.2 字节转换

	// 对象转字节数组
	bytes, err := XObject.ToByte(obj)

	// 字节数组转对象
	err := XObject.FromByte(bytes, &obj)

更多信息请参考模块文档。
*/
package XObject
