# XUtility

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XUtility.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XUtility)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)

XUtility 提供了常用的工具函数集，如数值计算、随机数等。

## 功能特性

- 数值计算：提供最大值、最小值计算
- 随机数生成：支持指定范围的随机整数生成

## 使用手册

### 1. 数值计算

#### 1.1 获取最大最小值

```go
// 获取两个数中的最大值
max := XUtility.MaxValue(10, 20) // 返回 20

// 获取两个数中的最小值
min := XUtility.MinValue(10, 20) // 返回 10
```

### 2. 随机数生成

#### 2.1 生成指定范围的随机数

```go
// 生成 [1, 100) 范围内的随机数
// 注意：右区间是开区间，不包含 100
rand := XUtility.RandInt(1, 100)
```

## 常见问题

### 1. 随机数范围说明

RandInt 函数生成的随机数范围是左闭右开区间 [min, max)，即包含最小值但不包含最大值。例如：
- RandInt(1, 10) 生成的随机数可能是：1, 2, 3, 4, 5, 6, 7, 8, 9
- 不会生成 10

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
