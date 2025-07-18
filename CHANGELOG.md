# 更新记录

## [0.0.7] - 2025-07-08
### 修复
- 修复 XApp.Run 函数在应用退出时未完整调用 defer 的问题

## [0.0.6] - 2025-07-05
### 修复
- 修复 XCollect.Map 不严格的锁引起的数据竞态问题

## [0.0.5] - 2025-05-30
### 修复
- 修复 XLoom 定时器（timer.go）模块回调频率及 panic 恢复逻辑错误
- 修复 XLoom 定时器（timer.go）模块错误的单元测试

## [0.0.4] - 2025-05-29
### 新增
- 新增 ${Env.NumCPU} 变量用于引用求值并完善了 XEnv.Eval 的单元测试

### 修复
- 修复 XLoom 线程未监听退出信号的问题

## [0.0.3] - 2025-05-24
### 新增
- 新增 XCollect.Map 模块，适用于高并发、读多于写的业务场景

### 变更
- 优化 XCollect.Map 模块的文档

## [0.0.2] - 2025-05-06
### 变更
- 修改 XLog 模块的日志标签格式为：[key1=value1, ...]
- 优化 XLoom 模块的单元测试

### 修复
- 修复 XLog 模块的日志标签 Clone 时未复制 level 字段的问题
- 修正 XLog 模块若干 lint 的警告

### 新增
- 新增 XLoom 模块的指标度量（Prometheus）

## [0.0.1] - 2025-03-23
### 新增
- 首次发布
