# XEnv

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XEnv.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XEnv)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-Explore-blue)](https://deepwiki.com/eframework-org/GO.UTIL)

XEnv 是一个环境配置管理工具，支持多平台识别、应用配置管理、路径管理、命令行参数解析和环境变量求值等功能。

## 功能特性

- 参数解析：支持多种参数形式和缓存管理
- 环境配置：支持应用类型、运行模式、版本等环境配置
- 变量求值：支持 ${Env.Key} 格式的环境变量引用和求值
- 路径管理：提供本地路径和资产路径的统一管理

## 使用手册

### 1. 命令行参数

#### 1.1 参数解析
```go
// 命令行：./app --config=dev.json --port 8080
config := XEnv.GetArg("config")    // 返回 "dev.json"
port := XEnv.GetArg("port")        // 返回 "8080"

// 获取所有参数
args := XEnv.GetArgs()             // 返回参数 map
```

### 2. 环境配置

#### 2.1 应用类型
```go
// 获取应用类型和运行模式
appType := XEnv.App()             // 返回 Server/Client
appMode := XEnv.Mode()            // 返回 Dev/Test/Staging/Prod
version := XEnv.Version()         // 返回应用版本
```

#### 2.2 项目信息
```go
// 获取项目相关信息
solution := XEnv.Solution()        // 返回解决方案名称
project := XEnv.Project()         // 返回项目名称
product := XEnv.Product()         // 返回产品名称
channel := XEnv.Channel()         // 返回渠道名称
```

### 3. 变量求值

#### 3.1 内置变量
```go
// 获取应用信息
app := XEnv.Vars().Eval("${Env.App}")              // 获取应用类型
mode := XEnv.Vars().Eval("${Env.Mode}")            // 获取运行模式
platform := XEnv.Vars().Eval("${Env.Platform}")    // 获取运行平台

// 获取路径信息
localPath := XEnv.Vars().Eval("${Env.LocalPath}")  // 获取本地路径
assetPath := XEnv.Vars().Eval("${Env.AssetPath}")  // 获取资产路径
```

#### 3.2 参数引用
```go
// 获取命令行参数
config := XEnv.Vars().Eval("${Env.config}")        // 获取 --config 的值

// 获取系统环境变量
path := XEnv.Vars().Eval("${Env.PATH}")           // 获取系统 PATH 变量
```

#### 3.3 特殊处理
```go
// 嵌套变量
value := XEnv.Vars().Eval("${Env.${Env.KEY}}")    // 返回 "(Nested)"

// 循环引用
value := XEnv.Vars().Eval("${Env.A${Env.B}}")     // 返回 "(Recursive)"

// 未知变量
value := XEnv.Vars().Eval("${Env.UNKNOWN}")       // 返回 "(Unknown)"
```

### 4. 路径管理

#### 4.1 本地路径
```go
// 获取本地数据路径（自动创建目录）
localPath := XEnv.LocalPath()     // 返回规范化的本地路径
```

#### 4.2 资产路径
```go
// 获取资产文件路径
assetPath := XEnv.AssetPath()     // 返回规范化的资产路径
```

## 常见问题

### 1. 命令行参数格式有什么要求？

命令行参数支持两种格式：
1. `--key=value` 格式：参数名和值通过等号连接
2. `--key value` 格式：参数名和值通过空格分隔

示例：
```bash
./app --config=dev.json --port 8080 --debug=true
```

### 2. 变量求值有哪些注意事项？

1. 变量引用格式：使用 `${Env.Key}` 格式
2. 变量优先级：内置变量 > 命令行参数 > 系统环境变量
3. 特殊情况：
   - 嵌套变量：返回 `(Nested)` 后缀
   - 循环引用：返回 `(Recursive)` 后缀
   - 未知变量：返回 `(Unknown)` 后缀

### 3. 路径管理有什么注意事项？

1. 本地数据路径：
   - 自动创建目录（权限：0755）
   - 默认路径：`Local`

2. 资产文件路径：
   - 默认路径：`Assets`
   - 不会自动创建目录

3. 通用特性：
   - 自动转换为绝对路径
   - 统一使用斜杠（/）作为分隔符
   - 支持命令行参数指定：
     ```bash
     ./app --LocalPath=/data/local --AssetPath=/data/assets
     ```

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE) 