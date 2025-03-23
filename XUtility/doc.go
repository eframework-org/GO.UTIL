// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XUtility 提供了常用的工具函数集，如数值计算、随机数等。

功能特性

  - 数值计算：提供最大值、最小值计算
  - 随机数生成：支持指定范围的随机整数生成

使用示例

1. 数值计算

1.1 获取最大最小值

	max := XUtility.MaxValue(10, 20) // 返回 20
	min := XUtility.MinValue(10, 20) // 返回 10

2. 随机数生成

2.1 生成指定范围的随机数

	// 生成 [1, 100) 范围内的随机数
	rand := XUtility.RandInt(1, 100)

更多信息请参考模块文档。
*/
package XUtility
