# XPrefs

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XPrefs.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XPrefs)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)

XPrefs 是一个灵活高效的配置系统，实现了多源化配置的读写，支持变量求值和命令行参数覆盖等功能。

## 功能特性

- 多源化配置：支持内置配置（只读）、本地配置（可写）和远程配置（只读），支持多个配置源按优先级顺序读取
- 多数据类型：支持基础类型（整数、浮点数、布尔值、字符串）、数组类型及配置实例（IBase）
- 变量求值：支持通过命令行参数动态覆盖配置项，使用 ${Prefs.Key} 语法引用其他配置项

## 使用手册

### 1. 基础配置操作

#### 1.1 获取配置实例

```go
// 获取资产配置（只读）
asset := XPrefs.Asset()

// 获取本地配置（读写）
local := XPrefs.Local()

// 创建新的配置实例
config := XPrefs.New()
```

#### 1.2 读写配置项

```go
// 检查配置项是否存在
exists := XPrefs.HasKey("key")

// 设置配置项
local.Set("key", "value")

// 获取配置项（支持默认值）
value := XPrefs.Get("key", "default")

// 删除配置项
local.Unset("key")
```

#### 1.3 多源配置

```go
// 1. 从指定的配置源中查找
value := XPrefs.Get("key", "default", local, asset)  // 优先从 local 查找，然后是 asset

// 2. 从默认配置源中查找（仅 asset）
value := XPrefs.Get("key", "default")  // 等同于 asset.Get("key")

// 3. 检查配置项是否存在于任意源
exists := XPrefs.HasKey("key", local, asset)  // 依次检查 local 和 asset

// 4. 获取不同类型的值
intVal := XPrefs.GetInt("number", 0, local, asset)      // 整数
floatVal := XPrefs.GetFloat("price", 0.0, local, asset) // 浮点数
strVal := XPrefs.GetString("name", "", local, asset)    // 字符串
boolVal := XPrefs.GetBool("flag", false, local, asset)  // 布尔值

// 5. 获取数组类型的值
intArray := XPrefs.GetInts("numbers", []int{}, local, asset)           // 整数数组
floatArray := XPrefs.GetFloats("prices", []float32{}, local, asset)    // 浮点数数组
strArray := XPrefs.GetStrings("names", []string{}, local, asset)       // 字符串数组
boolArray := XPrefs.GetBools("flags", []bool{}, local, asset)          // 布尔值数组
```

配置源优先级规则：
1. 按照传入的配置源顺序依次查找
2. 如果所有指定的配置源都未找到，则查找资产配置（asset）
3. 如果资产配置也未找到，则返回默认值
4. 如果未指定配置源且未提供默认值，则仅从资产配置中查找

### 2. 类型转换

#### 2.1 基础类型

```go
// 获取整数值
intVal := XPrefs.GetInt("number", 0)

// 获取浮点数值
floatVal := XPrefs.GetFloat("price", 0.0)

// 获取布尔值
boolVal := XPrefs.GetBool("enabled", false)

// 获取字符串值
strVal := XPrefs.GetString("name", "")
```

#### 2.2 数组类型

```go
// 获取整数数组
intArray := XPrefs.GetInts("numbers")

// 获取浮点数数组
floatArray := XPrefs.GetFloats("prices")

// 获取布尔值数组
boolArray := XPrefs.GetBools("flags")

// 获取字符串数组
strArray := XPrefs.GetStrings("names")
```

### 3. 多级配置

#### 3.1 设置多级配置

```go
// 创建并设置多级配置
local.Set("UI", XPrefs.New()).
    Get("UI").(XPrefs.IBase).
    Set("Window", XPrefs.New()).
    Get("Window").(XPrefs.IBase).
    Set("Style", XPrefs.New()).
    Get("Style").(XPrefs.IBase).
    Set("Theme", "Dark")

// 获取多级配置
theme := local.Get("UI").(XPrefs.IBase).
    Get("Window").(XPrefs.IBase).
    Get("Style").(XPrefs.IBase).
    GetString("Theme")  // 返回 "Dark"
```

#### 3.2 命令行参数设置

```bash
# 使用命令行参数设置多级配置
./program --Prefs@Local.UI.Window.Style.Theme=Dark
```

### 4. 变量引用

#### 4.1 基本引用

```go
// 设置配置项
local.Set("name", "John")
local.Set("greeting", "Hello ${Prefs.name}")

// 解析变量引用
message := local.Eval("${Prefs.greeting}")  // 返回 "Hello John"
```

#### 4.2 多级引用

```go
// 设置多级配置
local.Set("user", XPrefs.New()).
    Get("user").(XPrefs.IBase).
    Set("name", "John")

// 解析多级变量引用
message := local.Eval("User: ${Prefs.user.name}")  // 返回 "User: John"
```

### 5. 配置同步

#### 5.1 自动保存

```go
// 配置会在程序退出时自动保存
// 也可以手动调用保存
local.Save()
```

#### 5.2 信号监听

系统会自动监听以下信号并在接收到信号时保存配置：
- SIGTERM：终止信号
- SIGINT：中断信号（Ctrl+C）

## 常见问题

### 1. 配置文件在哪里？
默认情况下，配置文件位于程序运行目录：
- 资产配置：`Assets/Preferences.json`
- 本地配置：`Local/Preferences.json`

### 2. 如何指定配置文件路径？
可以通过命令行参数指定：
```bash
./program --Prefs@Asset=path/to/asset.json --Prefs@Local=path/to/local.json
```

### 3. 如何使用命令行参数？
命令行参数支持以下几种用法：
```bash
# 指定配置文件路径
./program --Prefs@Asset=path/to/asset.json --Prefs@Local=path/to/local.json

# 直接设置配置项值
./program --Prefs@Asset.name=value

# 设置多级配置项
./program --Prefs@Local.UI.Theme=Dark

# 设置日志级别
./program --Prefs.Log.Level=Debug
```

### 4. 变量引用的特殊情况？
- 循环引用：将显示为 `(Recursive)`
- 未定义变量：将显示为 `(Unknown)`
- 空值：将显示为 `(Unknown)`
- 嵌套变量：将显示为 `(Nested)`

### 5. 性能优化建议

1. 配置访问优化：
   - 缓存常用配置值
   - 避免频繁读取配置
   - 使用批量获取替代多次获取

2. 内存使用优化：
   - 及时清理不需要的配置项
   - 避免过深的配置层级
   - 合理使用数组类型配置

3. 文件操作优化：
   - 合理设置自动保存间隔
   - 避免频繁手动保存
   - 使用异步保存机制

4. 变量引用优化：
   - 减少变量嵌套层级
   - 避免复杂的变量引用链
   - 使用缓存优化求值性能

### 6. 最佳实践建议

1. 配置组织：
   - 按功能模块划分配置
   - 保持配置层级扁平
   - 使用有意义的配置键名

2. 错误处理：
   - 始终提供合理的默认值
   - 验证配置值的有效性
   - 记录配置错误信息

3. 安全性：
   - 敏感配置使用加密存储
   - 控制配置文件权限
   - 避免明文存储密码

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
