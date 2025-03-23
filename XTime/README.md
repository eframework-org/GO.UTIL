# XTime

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XTime.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XTime)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)

XTime 提供了一组时间处理工具函数，支持时间戳转换和格式化等功能。

## 功能特性

- 时间戳获取：支持微秒、毫秒和秒级精度
- 时间转换：在时间戳和time.Time对象间转换
- 零点时间：计算指定时间的零点和到零点的时间差
- 时间格式化：支持多种预定义格式模板

## 使用手册

### 1. 时间戳获取

#### 1.1 获取不同精度的时间戳
```go
// 获取微秒级时间戳
microsecond := XTime.GetMicrosecond()

// 获取毫秒级时间戳
millisecond := XTime.GetMillisecond()

// 获取秒级时间戳
timestamp := XTime.GetTimestamp()
```

### 2. 时间转换

#### 2.1 时间戳与time.Time转换
```go
// 获取当前时间对象
now := XTime.NowTime()

// 时间戳转time.Time对象
timestamp := XTime.GetTimestamp()
timeObj := XTime.ToTime(timestamp)
```

### 3. 零点时间计算

#### 3.1 零点相关计算
```go
// 获取当前时间到下一个零点的秒数
seconds := XTime.TimeToZero()

// 获取指定时间戳的零点时间戳
zeroTimestamp := XTime.ZeroTime(timestamp)
```

### 4. 时间格式化

#### 4.1 使用预定义格式
```go
timestamp := XTime.GetTimestamp()

// 使用标准格式 (2006-01-02 15:04:05 +0800 CST)
fullTime := XTime.Format(timestamp, XTime.FormatFull)

// 使用简易格式 (2006-01-02 15:04:05)
liteTime := XTime.Format(timestamp, XTime.FormatLite)

// 使用文件名格式 (2006-01-02_15_04_05)
fileTime := XTime.Format(timestamp, XTime.FormatFile)
```

### 5. 预定义常量

#### 5.1 时间格式模板
```go
FormatFull = "2006-01-02 15:04:05 +0800 CST" // 标准格式
FormatLite = "2006-01-02 15:04:05"           // 简易格式
FormatFile = "2006-01-02_15_04_05"           // 文件名格式
```

#### 5.2 时间常量（秒）
```go
Second1  = 1      // 1秒
Minute1  = 60     // 1分钟
Hour1    = 3600   // 1小时
Day1     = 86400  // 1天
```

## 常见问题

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
