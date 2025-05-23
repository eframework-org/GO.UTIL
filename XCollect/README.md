# XCollect

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XCollect.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XCollect)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-Explore-blue)](https://deepwiki.com/eframework-org/GO.UTIL)

XCollect 提供了一组集合类型数据的工具函数集，包括数组操作工具和线程安全的字典工具。

## 功能特性

- 数组工具：泛型数组操作函数，支持查找、删除、插入等常用功能，适用于任意可比较类型
- 字典工具：线程安全的泛型 Map，支持高效的读写操作和顺序/并发遍历，兼具性能与易用性

## 使用手册

### 1. 数组工具

#### 1.1 按值查找
使用 `Index` 和 `Contains` 函数进行精确值匹配：

```go
arr := []int{1, 2, 3, 4, 5}

// 查找元素索引
idx := XCollect.Index(arr, 3)        // 返回 2

// 检查元素是否存在
exists := XCollect.Contains(arr, 3)   // 返回 true
```

#### 1.2 条件查找
使用函数作为查找条件，实现灵活的查找逻辑：

```go
// 查找第一个偶数
idx := XCollect.Index(arr, func(x int) bool {
    return x%2 == 0
})                                    // 返回 1（元素2的索引）
```

#### 1.3 元素删除
提供两种删除方式：按值删除和按索引删除：

```go
// 按值删除（删除所有匹配的元素）
arr = XCollect.Remove(arr, 3)         // 返回 [1, 2, 4, 5]

// 按索引删除（删除指定位置的元素）
arr = XCollect.Delete(arr, 1)         // 返回 [1, 4, 5]
```

#### 1.4 元素添加
支持在数组末尾追加或在指定位置插入元素：

```go
// 在末尾追加元素
arr = XCollect.Append(arr, 6)         // 返回 [1, 4, 5, 6]

// 在指定位置插入元素
arr = XCollect.Insert(arr, 1, 2)      // 返回 [1, 2, 4, 5, 6]
```

### 2. 字典工具

#### 2.1 基本操作

创建和基本的增删改查操作：

```go
// 创建一个新的线程安全Map
map := XCollect.NewMap()

// 存储键值对
map.Store("key1", 100)
map.Store("key2", 200)

// 读取值
value, exists := map.Load("key1")  // 返回 100, true

// 删除键值对
map.Delete("key1")

// 清空所有键值对
map.Clear()
```

#### 2.2 遍历操作
两种遍历方式，适用于不同的场景：

```go
// 基础遍历
map.Range(func(key, value any) bool {
    fmt.Printf("键: %v, 值: %v\n", key, value)
    return true  // 返回false可以提前终止遍历
})

// 并发遍历（适用于大数据）
map.RangeConcurrent(func(chunk int, key, value any) bool {
    fmt.Printf("分片: %d, 键: %v, 值: %v\n", chunk, key, value)
    return true  // 返回false可以提前终止所有协程的遍历
}, func(chunk int) {
    fmt.Printf("开始并发遍历，分片数量: %d\n", chunk)
})
```

## 常见问题

### 1. XCollect 数组工具函数支持哪些数据类型？
支持所有满足 `comparable` 约束的类型，包括：
- 基本类型：整数、浮点数、字符串等
- 复合类型：结构体（需要可比较）、指针等
注意：
- 使用泛型避免接口转换开销
- 就地修改数组减少内存分配
- 使用 `append` 优化切片操作

### 2. XCollect.Map 的性能及适用场景？

#### 2.1 读写操作

📊 `XCollect.Map` vs `sync.Map` 性能对照表（数据量 `10000`）：

| Map 类型      | CPU 核数 | 操作次数 (N)  | 平均时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) |
| ------------ | ------ | --------- | ------------ | ----------- | ---------------- |
| **XCollect.Map** | 1      | 17799663  | 65.10 ns/op  | 23 B/op     | 2 allocs/op      |
| XCollect.Map | 2      | 23778950  | 48.38 ns/op  | 23 B/op     | 2 allocs/op      |
| XCollect.Map | 4      | 28436152  | 40.58 ns/op  | 23 B/op     | 2 allocs/op      |
| XCollect.Map | 8      | 36565410  | 35.37 ns/op  | 23 B/op     | 2 allocs/op      |
| XCollect.Map | 16     | 36156556  | 32.37 ns/op  | 23 B/op     | 2 allocs/op      |
| XCollect.Map | 32     | 37887278  | 32.13 ns/op  | 23 B/op     | 2 allocs/op      |
| **sync.Map**     | 1      | 19024426  | 56.36 ns/op  | 23 B/op     | 1 allocs/op      |
| sync.Map     | 2      | 31349271  | 37.84 ns/op  | 23 B/op     | 1 allocs/op      |
| sync.Map     | 4      | 56206615  | 21.92 ns/op  | 23 B/op     | 1 allocs/op      |
| sync.Map     | 8      | 91515728  | 14.89 ns/op  | 23 B/op     | 1 allocs/op      |
| sync.Map     | 16     | 100000000 | 11.17 ns/op  | 23 B/op     | 1 allocs/op      |
| sync.Map     | 32     | 132207163 | 9.795 ns/op  | 23 B/op     | 1 allocs/op      |

数据分析：

1. 两者都表现出良好的扩展性，`CPU` 核数越多，`ns/op` 越低，但 `sync.Map` 在多核下的性能提升更为显著，尤其是 `8 核` 及以上时表现优越。
2. `XCollect.Map` 尽管平均性能逊于 `sync.Map`，但在低并发场景下（`1–4 核`）仍表现出较强竞争力，平均时延接近或略优。
3. 两者单位操作的内存占用相同（`23 B/op`），`sync.Map` 每次操作只产生一次内存分配，而 `XCollect.Map` 每次产生两次分配，可能影响 `GC` 压力。

#### 2.2 遍历操作

📊 `XCollect.Map Range` vs `XCollect.Map Concurrent Range` vs `sync.Map Range` vs `map range` 性能对照表（数据量 `100000`）：

| Map 类型     | CPU 核数 | 操作次数 (N) | 平均时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) |
| ----------- | ------ | -------- | ------------ | ----------- | ---------------- |
| **XCollect.Map** | 1      | 10000    | 106382 ns/op  | 0 B/op      | 0 allocs/op      |
| XCollect.Map | 2      | 12199    | 104806 ns/op  | 0 B/op      | 0 allocs/op      |
| XCollect.Map | 4      | 10000    | 101121 ns/op  | 0 B/op      | 0 allocs/op      |
| XCollect.Map | 8      | 10000    | 102604 ns/op  | 0 B/op      | 0 allocs/op      |
| XCollect.Map | 16     | 10000    | 100402 ns/op  | 0 B/op      | 0 allocs/op      |
| XCollect.Map | 32     | 10000    | 105550 ns/op  | 0 B/op      | 0 allocs/op      |
| **XCollect.Map(Concurrent)** | 1      | 11443    | 110427 ns/op  | 2068 B/op   | 66 allocs/op     |
| XCollect.Map(Concurrent) | 2      | 19808    | 63128 ns/op   | 2068 B/op   | 66 allocs/op     |
| XCollect.Map(Concurrent) | 4      | 30472    | 42135 ns/op   | 2068 B/op   | 66 allocs/op     |
| XCollect.Map(Concurrent) | 8      | 34371    | 41602 ns/op   | 2068 B/op   | 66 allocs/op     |
| XCollect.Map(Concurrent) | 16     | 29734    | 42611 ns/op   | 2068 B/op   | 66 allocs/op     |
| XCollect.Map(Concurrent) | 32     | 27742    | 41233 ns/op   | 2068 B/op   | 66 allocs/op     |
| **sync.Map** | 1      | 584      | 2023889 ns/op | 0 B/op      | 0 allocs/op      |
| sync.Map     | 2      | 589      | 2018034 ns/op | 0 B/op      | 0 allocs/op      |
| sync.Map     | 4      | 588      | 2066907 ns/op | 0 B/op      | 0 allocs/op      |
| sync.Map     | 8      | 566      | 2082428 ns/op | 0 B/op      | 0 allocs/op      |
| sync.Map     | 16     | 572      | 2169874 ns/op | 0 B/op      | 0 allocs/op      |
| sync.Map     | 32     | 639      | 2033276 ns/op | 0 B/op      | 0 allocs/op      |
| **map**          | 1      | 2276     | 527344 ns/op  | 0 B/op      | 0 allocs/op      |
| map          | 2      | 2293     | 523724 ns/op  | 0 B/op      | 0 allocs/op      |
| map          | 4      | 2408     | 517579 ns/op  | 0 B/op      | 0 allocs/op      |
| map          | 8      | 2262     | 536811 ns/op  | 0 B/op      | 0 allocs/op      |
| map          | 16     | 2329     | 528598 ns/op  | 0 B/op      | 0 allocs/op      |
| map          | 32     | 2270     | 533023 ns/op  | 0 B/op      | 0 allocs/op      |

数据分析：

1. `XCollect.Map` 普通遍历在全核数范围内稳定在 `100μs` 左右，非常稳定且零分配。
2. `XCollect.Map` 并发遍历利用多核并发（`4–8` 核效率最优），可将遍历时间降到约 `41μs（~2.4 倍加速）`，但带来小量额外分配（每次遍历约 `2KB` 内存、`66` 次分配）。
3. `sync.Map` 遍历性能较差差，平均遍历耗时超过 `2ms（2000μs）`。
4. `map` 的遍历速度约在 `500μs` 水平，低于 `Sync.Map`，但相较于 `XCollect.Map` 没有优势且不支持并发。

#### 2.3 适用场景

1. `XCollect.Map`：高并发、大数据、读多写少。
2. `sync.Map`：高并发、数据可控、读写均衡。
3. `map`：低并发、数据可控、读写均衡。

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
