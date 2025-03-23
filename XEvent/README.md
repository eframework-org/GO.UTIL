# XEvent

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XEvent.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XEvent)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)

XEvent 是一个轻量级的事件管理器，支持多重监听、单次回调和批量通知等功能。

## 功能特性

- 多重监听：可配置是否允许同一事件注册多个回调
- 单次执行：可设置回调函数仅执行一次后自动注销

## 使用手册

### 1. 事件管理

#### 1.1 创建管理器
```go
// 创建一个支持多重注册的事件管理器
mgr := XEvent.NewManager(true)

// 获取共享的事件管理器
mgr := XEvent.Shared()
```

#### 1.2 注册事件
```go
// 注册事件处理程序
mgr.Reg(1001, func(args ...any) {
    // 处理事件
})

// 注册只执行一次的处理程序
mgr.Reg(1001, func(args ...any) {
    // 处理事件
}, true)
```

#### 1.3 注销事件
```go
// 注销指定的事件处理程序
mgr.Unreg(1001, handler)

// 注销事件的所有处理程序
mgr.Unreg(1001)
```

#### 1.4 通知事件
```go
// 通知事件（不带参数）
mgr.Notify(1001)

// 通知事件（带参数）
mgr.Notify(1001, "param1", 123)
```

## 常见问题

### 1. 事件处理程序的注册限制
1. 默认模式（Multiple=false）：
   - 每个事件 ID 只能注册一个处理程序
   - 重复注册会返回 false

2. 多重模式（Multiple=true）：
   - 每个事件 ID 可以注册多个处理程序
   - 按注册顺序依次执行

### 2. 事件通知的执行顺序
1. 同步执行：
   - 处理程序按注册顺序依次执行
   - 每个处理程序在当前 goroutine 中执行

2. 异常处理：
   - 自动捕获处理程序中的异常
   - 不影响其他处理程序的执行

### 3. 对象池的使用说明
1. 对象复用：
   - 事件包装器（EvtWrap）和处理程序包装器（HndWrap）使用对象池
   - 注销事件时自动回收对象到池中

2. 性能优化：
   - 减少内存分配和 GC 压力
   - 适合高频事件处理场景

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE) 