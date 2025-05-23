// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XCollect 提供了一组集合类型数据的工具函数集，包括数组操作工具和线程安全的字典工具。

功能特性

- 数组工具：泛型数组操作函数，支持查找、删除、插入等常用功能，适用于任意可比较类型

- 字典工具：线程安全的泛型 Map，支持高效的读写操作和顺序/并发遍历，兼具性能与易用性

使用手册

1. 数组工具

1.1 按值查找

	arr := []int{1, 2, 3, 4, 5}
	idx := XCollect.Index(arr, 3)        // 返回 2
	exists := XCollect.Contains(arr, 3)   // 返回 true

1.2 条件查找

	idx = XCollect.Index(arr, func(x int) bool {
	    return x%2 == 0
	})                                    // 返回第一个偶数的索引 1

1.3 元素删除

	arr = XCollect.Remove(arr, 3)         // 返回 [1, 2, 4, 5]
	arr = XCollect.Delete(arr, 1)         // 返回 [1, 4, 5]

1.4 元素添加

	arr = XCollect.Append(arr, 6)         // 返回 [1, 4, 5, 6]
	arr = XCollect.Insert(arr, 1, 2)      // 返回 [1, 2, 4, 5, 6]

2. 字典工具

2.1 基本操作

	map := XCollect.NewMap()
	map.Store("key1", 100)
	map.Store("key2", 200)
	value, exists := map.Load("key1")  // 返回 100, true
	map.Delete("key1")
	map.Clear()

2.2 遍历操作

	map.Range(func(key, value any) bool {
	    fmt.Printf("键: %v, 值: %v\n", key, value)
	    return true
	})

	map.RangeConcurrent(func(chunk int, key, value any) bool {
	    fmt.Printf("分片: %d, 键: %v, 值: %v\n", chunk, key, value)
	    return true
	}, func(chunk int) {
	    fmt.Printf("开始并发遍历，分片数量: %d\n", chunk)
	})

更多信息请参考模块文档。
*/
package XCollect
