// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/shiena/ansicolor"
)

// newStdBursh 创建一个用于给文本添加 ANSI 颜色的函数。
// 根据指定的颜色代码生成一个可以为文本添加颜色的格式化函数。
func newStdBursh(color string) func(string) string {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string { return pre + color + "m" + text + reset }
}

// stdBrushes 定义了不同日志级别对应的 ANSI 颜色格式化函数。
// Emergency 使用黑色，Alert 使用青色，Critical 使用品红色，
// Error 使用红色，Warn 使用黄色，Notice 使用绿色，
// Info 使用灰色，Debug 使用蓝色。
var stdBrushes = []func(string) string{
	newStdBursh("1;39"), // Emergency          black
	newStdBursh("1;36"), // Alert              cyan
	newStdBursh("1;35"), // Critical           magenta
	newStdBursh("1;31"), // Error              red
	newStdBursh("1;33"), // Warn               yellow
	newStdBursh("1;32"), // Notice             green
	newStdBursh("1;30"), // Info               grey
	newStdBursh("1;34"), // Debug              blue
}

// 标准输出适配器的配置项常量
const (
	stdPrefsLevel        = "Level"      // 日志级别配置项名称
	stdPrefsLevelDefault = LevelInfoStr // 默认日志级别
	stdPrefsColor        = "Color"      // 颜色开关配置项名称
	stdPrefsColorDefault = true         // 默认启用颜色输出
)

// stdAdapter 实现了标准输出日志适配器。
// 支持日志级别过滤和 ANSI 颜色输出。
type stdAdapter struct {
	level  LevelType // 当前日志级别
	color  bool      // 是否启用颜色输出
	writer io.Writer // 输出目标
}

// newStdAdapter 创建一个新的标准输出日志适配器。
// 默认使用支持 ANSI 颜色的标准输出作为写入目标。
func newStdAdapter() *stdAdapter {
	apt := &stdAdapter{writer: ansicolor.NewAnsiColorWriter(os.Stdout)}
	return apt
}

// init 初始化标准输出日志适配器。
// 从配置中读取日志级别和颜色输出设置，并返回配置的日志级别。
func (apt *stdAdapter) init(prefs XPrefs.IBase) LevelType {
	if prefs == nil {
		return LevelUndefined
	}
	tmpLevel := prefs.GetString(stdPrefsLevel, stdPrefsLevelDefault)
	switch tmpLevel {
	case LevelDebugStr:
		apt.level = LevelDebug
	case LevelInfoStr:
		apt.level = LevelInfo
	case LevelNoticeStr:
		apt.level = LevelNotice
	case LevelWarnStr:
		apt.level = LevelWarn
	case LevelErrorStr:
		apt.level = LevelError
	case LevelCriticalStr:
		apt.level = LevelCritical
	case LevelAlertStr:
		apt.level = LevelAlert
	case LevelEmergencyStr:
		apt.level = LevelEmergency
	default:
		apt.level = LevelUndefined
	}
	apt.color = prefs.GetBool(stdPrefsColor, stdPrefsColorDefault)
	return apt.level
}

// write 将日志写入标准输出。
// 根据日志级别和颜色设置格式化日志内容，并写入到输出目标。
// 当日志为空时返回错误，当日志级别高于设定且未强制输出时跳过。
func (apt *stdAdapter) write(log *logData) error {
	if log == nil {
		return errors.New("nil log")
	}
	if log.level > apt.level && !log.force {
		return nil
	}
	str := log.text(true)
	if apt.color {
		str = strings.Replace(str, levelLabel[log.level], stdBrushes[log.level](levelLabel[log.level]), 1)
	}
	h, _, _ := formatTime(log.time)
	apt.writer.Write(append(append(h, str...), '\n'))
	return nil
}

// flush 刷新标准输出缓冲区。
// 标准输出适配器不需要特殊的刷新操作。
func (apt *stdAdapter) flush() {}

// close 关闭标准输出适配器。
// 标准输出适配器不需要特殊的关闭操作。
func (apt *stdAdapter) close() {}
