# XLog

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XLog.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XLog)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-Explore-blue)](https://deepwiki.com/eframework-org/GO.UTIL)

XLog 提供了一个遵循 RFC5424 标准的日志系统，支持多级别日志输出、多适配器管理、日志轮转和结构化标签等特性。

## 功能特性

- 支持 RFC5424 标准的 8 个日志级别
- 支持标准输出和文件存储两种适配器
- 支持日志文件的自动轮转和清理
- 支持异步写入和线程安全操作
- 支持结构化的日志标签系统

## 使用手册

### 1. 基础日志记录

#### 1.1 日志级别
```go
// 不同级别的日志记录
XLog.Emergency("系统崩溃")  // 级别 0：紧急
XLog.Alert("需要立即处理")   // 级别 1：警报
XLog.Critical("严重错误")   // 级别 2：严重
XLog.Error("操作失败")      // 级别 3：错误
XLog.Warn("潜在问题")       // 级别 4：警告
XLog.Notice("重要信息")     // 级别 5：通知
XLog.Info("一般信息")       // 级别 6：信息
XLog.Debug("调试信息")      // 级别 7：调试
```

### 2. 日志配置

#### 2.1 文件输出配置

文件输出适配器支持以下配置项：

```go
prefs := XPrefs.New()
fileConf := XPrefs.New()

// 基础配置
fileConf.Set("Path", "./logs/app.log")     // 日志文件路径，支持环境变量 ${Env.xxx}
fileConf.Set("Level", "Debug")             // 日志级别：Emergency|Alert|Critical|Error|Warn|Notice|Info|Debug

// 轮转配置
fileConf.Set("Rotate", true)               // 是否启用日志轮转，默认 true
fileConf.Set("Daily", true)                // 是否按天轮转，默认 true
fileConf.Set("MaxDay", 7)                  // 日志文件保留天数，默认 7 天
fileConf.Set("Hourly", true)               // 是否按小时轮转，默认 true
fileConf.Set("MaxHour", 168)               // 日志文件保留小时数，默认 168 小时（7天）

// 文件限制
fileConf.Set("MaxFile", 100)               // 最大文件数量，默认 100 个
fileConf.Set("MaxLine", 1000000)           // 单文件最大行数，默认 100 万行
fileConf.Set("MaxSize", 134217728)         // 单文件最大体积，默认 128MB

prefs.Set("Log/File", fileConf)
```

#### 2.2 标准输出配置

标准输出适配器支持以下配置项：

```go
stdConf := XPrefs.New()

// 基础配置
stdConf.Set("Level", "Info")               // 日志级别：Emergency|Alert|Critical|Error|Warn|Notice|Info|Debug
stdConf.Set("Color", true)                 // 是否启用彩色输出，默认 true

prefs.Set("Log/Std", stdConf)
```

#### 2.3 配置说明

1. 日志级别控制：
   - 通过配置每个适配器的 Level 参数控制，低于该级别的日志不会输出
   - 可以通过 XLog.Level() 获取当前最大级别
   - 可以通过 LogTag 的 Level() 方法设置特定标签的日志级别，这将优先于全局级别
   - 示例：
     ```go
     tag := XLog.GetTag()
     tag.Level(XLog.LevelDebug)  // 即使全局级别是 Info，带此标签的日志也会输出 Debug 级别
     XLog.Debug(tag, "调试信息")   // 此日志会被输出
     XLog.Watch(tag)
     XLog.Debug("调试信息")        // 同样会被输出，因为继承了上下文标签的级别
     XLog.Defer()                // 清除上下文标签
     ```

2. 日志级别优先级（从高到低）：
   - Emergency (0): 系统不可用
   - Alert (1): 必须立即采取措施
   - Critical (2): 严重错误
   - Error (3): 错误
   - Warn (4): 警告
   - Notice (5): 重要信息
   - Info (6): 一般信息
   - Debug (7): 调试信息

3. 文件轮转策略：
   - 按天轮转：每天创建新文件，自动清理超过 MaxDay 天数的文件
   - 按小时轮转：每小时创建新文件，自动清理超过 MaxHour 小时数的文件
   - 按大小轮转：当文件超过 MaxSize 时创建新文件
   - 按行数轮转：当文件超过 MaxLine 时创建新文件
   - 文件数量限制：通过 MaxFile 控制最大文件数

4. 日志文件命名：
   假设配置 Path 为 "./logs/app.log"：
   - 按天轮转：
     - 当前文件：app.log
     - 历史文件：app.2006-01-02.001.log, app.2006-01-02.002.log, ...
   - 按小时轮转：
     - 当前文件：app.log
     - 历史文件：app.2006-01-02-15.001.log, app.2006-01-02-15.002.log, ...
   - 按大小/行数轮转：
     - 当前文件：app.log
     - 历史文件：app.001.log, app.002.log, ...

   注意：
   - 如果 Path 配置中只有后缀（如 ".log"），则文件名部分为空，例如：
     - 当前文件：.log
     - 历史文件：2006-01-02.001.log（按天）, 2006-01-02-15.001.log（按小时）
   - 如果 Path 配置中只有目录，则使用空文件名和 ".log" 后缀
   - 历史文件的序号从 001 开始递增，最大受 MaxFile 参数限制
   - 日期格式使用 ISO 8601 标准（2006-01-02 表示年月日，15 表示小时）

5. 标准输出特性：
   - 支持 ANSI 颜色输出，不同日志级别使用不同颜色
   - Emergency: 黑色
   - Alert: 青色
   - Critical: 品红色
   - Error: 红色
   - Warn: 黄色
   - Notice: 绿色
   - Info: 灰色
   - Debug: 蓝色

### 3. 日志标签

#### 3.1 使用标签
```go
// 创建标签
tag := XLog.GetTag()
tag.Set("module", "auth")
tag.Set("user", "admin")

// 记录带标签的日志
XLog.Info(tag, "用户登录成功")

// 使用上下文标签
XLog.Watch(tag)
XLog.Info("执行操作") // 自动带上标签
XLog.Defer() // 清除上下文标签
```

### 4. 错误处理

#### 4.1 异常捕获
```go
defer XLog.Caught(false)
// 你的代码...
panic("发生错误")
```

## 常见问题

### 1. 日志文件没有轮转？
确保正确配置了轮转参数：
- `Rotate`: 是否启用轮转
- `Daily`/`Hourly`: 按天/小时轮转
- `MaxDay`/`MaxHour`: 文件保留时间
- `MaxLine`/`MaxSize`: 单文件限制

### 2. 日志级别如何控制？
通过配置每个适配器的 Level 参数控制，低于该级别的日志不会输出。可以通过 XLog.Level() 获取当前最大级别。此外，可以通过 LogTag 的 Level() 方法设置特定标签的日志级别，这将优先于全局级别。

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
