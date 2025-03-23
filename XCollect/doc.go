// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XCollect 提供了一组集合类型数据的工具函数集，支持泛型和函数查询。

功能特性

  - 数组操作：查找、插入、删除、追加等基础操作
  - 泛型支持：支持任意可比较类型的数据操作
  - 函数查询：支持使用函数作为查找和过滤条件
  - 类型安全：编译时类型检查，避免运行时错误

使用手册

1. 数组查找

1.1 按值查找

	arr := []int{1, 2, 3, 4, 5}
	idx := XCollect.Index(arr, 3)        // 返回 2
	exists := XCollect.Contains(arr, 3)   // 返回 true

1.2 按条件查找

	idx = XCollect.Index(arr, func(x int) bool {
	    return x%2 == 0
	})                                    // 返回第一个偶数的索引 1

2. 数组修改

2.1 元素删除

	arr = XCollect.Remove(arr, 3)         // 返回 [1, 2, 4, 5]
	arr = XCollect.Delete(arr, 1)         // 返回 [1, 4, 5]

2.2 元素添加

	arr = XCollect.Append(arr, 6)         // 返回 [1, 4, 5, 6]
	arr = XCollect.Insert(arr, 1, 2)      // 返回 [1, 2, 4, 5, 6]

更多信息请参考模块文档。
*/
package XCollect
