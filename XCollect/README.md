# XCollect

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XCollect.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XCollect)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-Explore-blue)](https://deepwiki.com/eframework-org/GO.UTIL)

XCollect 提供了一组集合类型数据的工具函数集，支持泛型和函数查询。

## 功能特性

- 数组操作：查找、插入、删除、追加等基础操作
- 泛型支持：支持任意可比较类型的数据操作
- 函数查询：支持使用函数作为查找和过滤条件
- 类型安全：编译时类型检查，避免运行时错误

## 使用手册

### 1. 数组查找

#### 1.1 按值查找
使用 `Index` 和 `Contains` 函数进行精确值匹配：

```go
arr := []int{1, 2, 3, 4, 5}

// 查找元素索引
idx := XCollect.Index(arr, 3)        // 返回 2

// 检查元素是否存在
exists := XCollect.Contains(arr, 3)   // 返回 true
```

#### 1.2 按条件查找
使用函数作为查找条件，实现灵活的查找逻辑：

```go
// 查找第一个偶数
idx := XCollect.Index(arr, func(x int) bool {
    return x%2 == 0
})                                    // 返回 1（元素2的索引）
```

### 2. 数组修改

#### 2.1 元素删除
提供两种删除方式：按值删除和按索引删除：

```go
// 按值删除（删除所有匹配的元素）
arr = XCollect.Remove(arr, 3)         // 返回 [1, 2, 4, 5]

// 按索引删除（删除指定位置的元素）
arr = XCollect.Delete(arr, 1)         // 返回 [1, 4, 5]
```

#### 2.2 元素添加
支持在数组末尾追加或在指定位置插入元素：

```go
// 在末尾追加元素
arr = XCollect.Append(arr, 6)         // 返回 [1, 4, 5, 6]

// 在指定位置插入元素
arr = XCollect.Insert(arr, 1, 2)      // 返回 [1, 2, 4, 5, 6]
```

## 常见问题

### 1. 如何处理空数组？
所有操作都会安全处理空数组：
- `Index` 和 `Contains` 返回 -1 和 false
- `Remove` 和 `Delete` 返回空数组
- `Append` 和 `Insert` 正常工作

### 2. 支持哪些数据类型？
支持所有满足 `comparable` 约束的类型，包括：
- 基本类型：整数、浮点数、字符串等
- 复合类型：结构体（需要可比较）、指针等

### 3. 性能如何优化？
- 使用泛型避免接口转换开销
- 就地修改数组减少内存分配
- 使用 `append` 优化切片操作

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
