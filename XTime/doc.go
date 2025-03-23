// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XTime 提供了一组时间处理工具函数，支持时间戳转换和格式化等功能。

功能特性

  - 时间戳获取：支持微秒、毫秒和秒级精度
  - 时间转换：在时间戳和time.Time对象间转换
  - 零点时间：计算指定时间的零点和到零点的时间差
  - 时间格式化：支持多种预定义格式模板

使用手册

1. 时间戳获取

1.1 获取不同精度的时间戳

	microsecond := XTime.GetMicrosecond() // 获取微秒级时间戳
	millisecond := XTime.GetMillisecond() // 获取毫秒级时间戳
	timestamp := XTime.GetTimestamp()      // 获取秒级时间戳

2. 时间转换

2.1 时间戳与time.Time转换

	now := XTime.NowTime()                // 获取当前时间对象
	timeObj := XTime.ToTime(timestamp)    // 时间戳转time.Time对象

3. 零点时间计算

3.1 零点相关计算

	seconds := XTime.TimeToZero()         // 获取当前时间到下一个零点的秒数
	zeroTimestamp := XTime.ZeroTime()     // 获取当前时间的零点时间戳

4. 时间格式化

4.1 使用预定义格式

	fullTime := XTime.Format(timestamp, XTime.FormatFull) // 标准格式
	liteTime := XTime.Format(timestamp, XTime.FormatLite) // 简易格式
	fileTime := XTime.Format(timestamp, XTime.FormatFile) // 文件名格式

更多信息请参考模块文档。
*/
package XTime
