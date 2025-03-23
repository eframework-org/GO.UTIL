# XFile

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XFile.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XFile)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)

XFile 提供了一组实用的文件系统操作工具函数，支持文件读写和路径处理等功能。

## 功能特性

- 文件操作：支持文件的常用操作，如：检查、读取、写入、删除等
- 目录操作：支持目录的常用操作，如：检查、创建、获取目录名称等
- 路径处理：统一的 POSIX 风格路径处理，支持路径拼接及标准化

## 使用手册

### 1. 文件操作

#### 1.1 文件检查
```go
// 检查文件是否存在
exists := XFile.HasFile("path/to/file.txt")
```

#### 1.2 文件读写
```go
// 读取二进制文件
data := XFile.OpenFile("config.dat")

// 写入二进制文件
err := XFile.SaveFile("config.dat", data, 0644)

// 读取文本文件
text := XFile.OpenText("config.txt")

// 写入文本文件
err := XFile.SaveText("config.txt", "Hello World", 0644)
```

### 2. 目录操作

#### 2.1 目录管理
```go
// 检查目录是否存在
exists := XFile.HasDirectory("path/to/dir")

// 创建目录
success := XFile.CreateDirectory("path/to/dir")

// 获取父目录
parent := XFile.DirectoryName("path/to/file.txt")
```

### 3. 路径处理

#### 3.1 路径操作
```go
// 连接路径
path := XFile.PathJoin("path", "to", "file.txt")

// 标准化路径
norm := XFile.NormalizePath("path/./to/../file.txt")
```

## 常见问题

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE) 