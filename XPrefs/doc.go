// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XPrefs 是一个灵活高效的配置系统，实现了多源化配置的读写，支持变量求值和命令行参数覆盖等功能。

功能特性

  - 多源化配置：支持内置配置（只读）、本地配置（可写）和远程配置（只读），支持多个配置源按优先级顺序读取
  - 多数据类型：支持基础类型（整数、浮点数、布尔值、字符串）、数组类型及配置实例（IBase）
  - 变量求值：支持通过命令行参数动态覆盖配置项，使用 ${Prefs.Key} 语法引用其他配置项

使用手册

1. 基础配置操作

1.1 获取配置实例

	// 获取资产配置（只读）
	asset := XPrefs.Asset()

	// 获取本地配置（读写）
	local := XPrefs.Local()

	// 创建新的配置实例
	config := XPrefs.New()

1.2 读写配置项

	// 检查配置项是否存在
	exists := XPrefs.HasKey("key")

	// 设置配置项
	local.Set("key", "value")

	// 获取配置项（支持默认值）
	value := XPrefs.Get("key", "default")

	// 删除配置项
	local.Unset("key")

1.3 多源配置

	// 从指定的配置源中查找
	value := XPrefs.Get("key", "default", local, asset)  // 优先从 local 查找，然后是 asset

	// 从默认配置源中查找（仅 asset）
	value := XPrefs.Get("key", "default")  // 等同于 asset.Get("key")

	// 检查配置项是否存在于任意源
	exists := XPrefs.HasKey("key", local, asset)  // 依次检查 local 和 asset

	// 获取不同类型的值
	intVal := XPrefs.GetInt("number", 0, local, asset)      // 整数
	floatVal := XPrefs.GetFloat("price", 0.0, local, asset) // 浮点数
	strVal := XPrefs.GetString("name", "", local, asset)    // 字符串
	boolVal := XPrefs.GetBool("flag", false, local, asset)  // 布尔值

	// 获取数组类型的值
	intArray := XPrefs.GetInts("numbers", []int{}, local, asset)           // 整数数组
	floatArray := XPrefs.GetFloats("prices", []float32{}, local, asset)    // 浮点数数组
	strArray := XPrefs.GetStrings("names", []string{}, local, asset)       // 字符串数组
	boolArray := XPrefs.GetBools("flags", []bool{}, local, asset)          // 布尔值数组

配置源优先级规则：
 1. 按照传入的配置源顺序依次查找
 2. 如果所有指定的配置源都未找到，则查找资产配置（asset）
 3. 如果资产配置也未找到，则返回默认值
 4. 如果未指定配置源且未提供默认值，则仅从资产配置中查找

2. 类型转换

2.1 基础类型

	// 获取整数值
	intVal := XPrefs.GetInt("number", 0)

	// 获取浮点数值
	floatVal := XPrefs.GetFloat("price", 0.0)

	// 获取布尔值
	boolVal := XPrefs.GetBool("enabled", false)

	// 获取字符串值
	strVal := XPrefs.GetString("name", "")

2.2 数组类型

	// 获取整数数组
	intArray := XPrefs.GetInts("numbers")

	// 获取浮点数数组
	floatArray := XPrefs.GetFloats("prices")

	// 获取布尔值数组
	boolArray := XPrefs.GetBools("flags")

	// 获取字符串数组
	strArray := XPrefs.GetStrings("names")

3. 多级配置

3.1 设置多级配置

	// 创建并设置多级配置
	local.Set("UI", XPrefs.New()).
		Get("UI").(XPrefs.IBase).
		Set("Window", XPrefs.New()).
		Get("Window").(XPrefs.IBase).
		Set("Style", XPrefs.New()).
		Get("Style").(XPrefs.IBase).
		Set("Theme", "Dark")

	// 获取多级配置
	theme := local.Get("UI").(XPrefs.IBase).
		Get("Window").(XPrefs.IBase).
		Get("Style").(XPrefs.IBase).
		GetString("Theme")  // 返回 "Dark"

3.2 命令行参数设置

命令行参数支持以下几种用法：
  - 指定配置文件路径：--Prefs@Asset=path/to/asset.json
  - 直接设置配置项值：--Prefs@Asset.name=value
  - 设置多级配置项：--Prefs@Local.UI.Theme=Dark
  - 设置日志级别：--Prefs.Log.Level=Debug

4. 变量引用

4.1 基本引用

	// 设置配置项
	local.Set("name", "John")
	local.Set("greeting", "Hello ${Prefs.name}")

	// 解析变量引用
	message := local.Eval("${Prefs.greeting}")  // 返回 "Hello John"

4.2 多级引用

	// 设置多级配置
	local.Set("user", XPrefs.New()).
		Get("user").(XPrefs.IBase).
		Set("name", "John")

	// 解析多级变量引用
	message := local.Eval("User: ${Prefs.user.name}")  // 返回 "User: John"

变量引用的特殊情况：
  - 循环引用：将显示为 (Recursive)
  - 未定义变量：将显示为 (Unknown)
  - 空值：将显示为 (Unknown)
  - 嵌套变量：将显示为 (Nested)

5. 配置同步

5.1 自动保存

	// 配置会在程序退出时自动保存
	// 也可以手动调用保存
	local.Save()

5.2 信号监听

系统会自动监听以下信号并在接收到信号时保存配置：
  - SIGTERM：终止信号
  - SIGINT：中断信号（Ctrl+C）

更多信息请参考模块文档。
*/
package XPrefs
