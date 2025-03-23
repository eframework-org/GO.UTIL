// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/eframework-org/GO.UTIL/XEnv"
	"github.com/eframework-org/GO.UTIL/XFile"
	"github.com/eframework-org/GO.UTIL/XTime"
	"github.com/illumitacit/gostd/quit"
)

// 时间格式化用的数字映射常量
const (
	y1  = `0123456789`                                                                                           // 年份数字1-9
	y2  = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789` // 年份数字0-9重复
	y3  = `0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999` // 月份数字映射
	y4  = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789` // 年份数字0-9重复
	mo1 = `000000000111`                                                                                         // 月份数字1-9
	mo2 = `123456789012`                                                                                         // 月份数字0-9
	d1  = `0000000001111111111222222222233`                                                                      // 日期数字1-9
	d2  = `1234567890123456789012345678901`                                                                      // 日期数字0-9
	h1  = `000000000011111111112222`                                                                             // 小时数字1-9
	h2  = `012345678901234567890123`                                                                             // 小时数字0-9
	mi1 = `000000000011111111112222222222333333333344444444445555555555`                                         // 分钟数字1-9
	mi2 = `012345678901234567890123456789012345678901234567890123456789`                                         // 分钟数字0-9
	s1  = `000000000011111111112222222222333333333344444444445555555555`                                         // 秒数字1-9
	s2  = `012345678901234567890123456789012345678901234567890123456789`                                         // 秒数字0-9
	ns1 = `0123456789`                                                                                           // 纳秒数字0-9
)

const unknownSource = "[?]" // 未知调用源标记

// 用于地址掩码的正则表达式和替换模式
var (
	ipPattern     = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)                   // IP地址匹配模式
	domainPattern = regexp.MustCompile(`([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}`) // 域名匹配模式
	portPattern   = regexp.MustCompile(`:\d+`)                                          // 端口号匹配模式
	ipMask        = []byte("**.**.**.**")                                               // IP地址掩码
	domainMask    = []byte("***.***")                                                   // 域名掩码
)

// formatTime 格式化时间戳为日志时间格式。
// time 为要格式化的时间。
// 返回格式化后的时间字节切片、日期和小时。
func formatTime(time time.Time) ([]byte, int, int) {
	_, mo, d := time.Date()
	h, mi, s := time.Clock()
	ns := time.Nanosecond() / 1000000
	var buf [21]byte

	buf[0] = '['
	buf[1] = mo1[mo-1]
	buf[2] = mo2[mo-1]
	buf[3] = '/'
	buf[4] = d1[d-1]
	buf[5] = d2[d-1]
	buf[6] = ' '
	buf[7] = h1[h]
	buf[8] = h2[h]
	buf[9] = ':'
	buf[10] = mi1[mi]
	buf[11] = mi2[mi]
	buf[12] = ':'
	buf[13] = s1[s]
	buf[14] = s2[s]
	buf[15] = '.'
	buf[16] = ns1[ns/100]
	buf[17] = ns1[ns%100/10]
	buf[18] = ns1[ns%10]
	buf[19] = ']'
	buf[20] = ' '

	return buf[0:], d, h
}

// formatLog 格式化日志内容。
// data 为日志数据，可以是字符串或其他类型。
// args 为可选的格式化参数。
// 返回格式化后的日志字符串。
func formatLog(data any, args ...any) string {
	var str string
	switch v := data.(type) {
	case string:
		str = v
		if len(args) == 0 {
			return str
		}
		if strings.Contains(str, "%") && !strings.Contains(str, "%%") {
			// format string
		} else {
			// do not contain format char
			str += strings.Repeat(" %v", len(args))
		}
	default:
		str = fmt.Sprint(v)
		if len(args) == 0 {
			return str
		}
		str += strings.Repeat(" %v", len(args))
	}
	return fmt.Sprintf(str, args...)
}

// UnAddr 对字符串中的IP地址和域名进行掩码处理。
// 保留端口号信息，将IP地址和域名替换为掩码。
// input 为包含IP地址和域名的字符串。
// 返回处理后的字符串。
func UnAddr(input string) string {
	inputBytes := []byte(input)

	// 替换IP地址
	inputBytes = ipPattern.ReplaceAllFunc(inputBytes, func(match []byte) []byte {
		return ipMask
	})

	// 替换域名
	inputBytes = domainPattern.ReplaceAllFunc(inputBytes, func(match []byte) []byte {
		return domainMask
	})

	// 保留端口
	inputBytes = portPattern.ReplaceAllFunc(inputBytes, func(match []byte) []byte {
		return match
	})

	return string(inputBytes)
}

// Caller 获取调用栈信息。
// stack 为要跳过的调用层级数。
// fullpath 为是否显示完整的文件路径，false 则只显示函数名。
// 返回格式化的调用栈信息。
func Caller(stack int, fullpath bool) string {
	if pc, file, line, ok := runtime.Caller(stack + 1); ok {
		if fullpath {
			return fmt.Sprintf("[%s:%d (0x%v)]", file, line, pc)
		} else {
			return fmt.Sprintf("[%s:%d (0x%v)]", runtime.FuncForPC(pc).Name(), line, pc)
		}
	}
	return unknownSource
}

// Trace 获取完整的错误堆栈信息。
// stack 为要跳过的调用层级数。
// err 为触发堆栈跟踪的错误。
// 返回格式化的堆栈信息和堆栈深度。
func Trace(stack int, err any) (string, int) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	start := stack + 1
	count := stack
	fmt.Fprintf(buf, "    skip %v\n", stack)
	for i := start; ; i++ {
		line := Caller(i, true)
		if line == unknownSource {
			break
		}
		count++
		fmt.Fprintf(buf, "    %v\n", line)
	}
	return buf.String(), count
}

// Elapse 创建一个用于计算函数执行时间的闭包。
// stack 为要跳过的调用层级数。
// callback 为可选的回调函数，在计时结束时调用。
// 返回一个在调用时输出执行时间的函数。
func Elapse(stack int, callback ...func()) func() {
	start := XTime.GetMillisecond()
	return func() {
		end := XTime.GetMillisecond()
		elapse := end - start
		if stack < 0 {
			stack = 0
		}
		caller := Caller(stack+1, false)
		Notice("XLog.Elapse%v: start time: %v, finish time: %v, elapsed-%vms", caller, start, end, elapse)
		if len(callback) == 1 {
			callback[0]()
		}
	}
}

// Caught 捕获并处理panic。
// exit 为是否在处理后退出程序。
// handler 为可选的自定义处理函数，接收错误信息和堆栈深度。
func Caught(exit bool, handler ...func(string, int)) {
	if err := recover(); err != nil {
		str, count := Trace(2, err) // 固定堆栈深度2
		fname := XFile.PathJoin(XEnv.LocalPath(), "Panic", fmt.Sprintf("%v.panic", XTime.Format(XTime.GetTimestamp(), XTime.FormatFile)))
		XFile.HasDirectory(XFile.DirectoryName(fname), true)
		XFile.SaveText(fname, str)
		Critical(str)
		if len(handler) == 1 {
			handler[0](str, count)
		}
		if exit {
			Critical("XLog.Caught: exit caused by panic.")
			quit.BroadcastShutdown()
			quit.GetWaiter().Wait()
			os.Exit(1)
		}
	}
}
