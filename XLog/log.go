// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/illumitacit/gostd/quit"
)

// LevelType 定义日志级别类型。
// 遵循 RFC5424 日志标准，定义了八个日志级别（0-7），用于表示日志消息的严重程度。
// 每个级别都有其特定的使用场景，从系统不可用的紧急情况到详细的调试信息。
type LevelType int

const (
	// LevelUndefined 表示未设置日志级别。
	// 用于初始化状态或表示无效的日志级别。
	LevelUndefined    LevelType = -1
	LevelUndefinedStr           = "Undefined"

	// LevelEmergency 表示最高级别的紧急情况（级别 0）。
	// 用于记录导致系统完全不可用的灾难性故障，需要立即通知所有技术人员。
	LevelEmergency    LevelType = 0
	LevelEmergencyStr           = "Emergency"

	// LevelAlert 表示需要立即处理的警报情况（级别 1）。
	// 用于记录需要立即采取行动的关键情况，如严重的安全事件或系统资源耗尽。
	LevelAlert    = 1
	LevelAlertStr = "Alert"

	// LevelCritical 表示严重的错误情况（级别 2）。
	// 用于记录严重的系统问题，如主要组件失效或关键功能不可用。
	LevelCritical    = 2
	LevelCriticalStr = "Critical"

	// LevelError 表示一般错误情况（级别 3）。
	// 用于记录影响系统正常运行的错误，如接口调用失败或业务流程中断。
	LevelError    = 3
	LevelErrorStr = "Error"

	// LevelWarn 表示警告情况（级别 4）。
	// 用于记录可能导致问题的异常情况，如资源即将耗尽或配置不当。
	LevelWarn    = 4
	LevelWarnStr = "Warn"

	// LevelNotice 表示重要的正常情况（级别 5）。
	// 用于记录需要注意但不属于错误的重要系统事件，如服务启动或配置更改。
	LevelNotice    = 5
	LevelNoticeStr = "Notice"

	// LevelInfo 表示一般信息（级别 6）。
	// 用于记录系统的正常运行信息，如常规操作日志或状态更新。
	LevelInfo    = 6
	LevelInfoStr = "Info"

	// LevelDebug 表示调试信息（级别 7）。
	// 用于记录详细的调试信息，帮助开发人员诊断和排查问题。
	LevelDebug    = 7
	LevelDebugStr = "Debug"
)

// levelMax 是可以输出的最大日志级别。
var levelMax LevelType

// logPool 是用于重用 logData 对象的 sync pool。
var logPool = sync.Pool{New: func() any { return &logData{} }}

// logCache 是用于缓存日志记录的通道。
var logCache = make(chan *logData, 300000)

// levelLabel 包含日志级别的字符串表示。
var levelLabel = [LevelDebug + 1]string{"[M]", "[A]", "[C]", "[E]", "[W]", "[N]", "[I]", "[D]"}

var (
	initMu    sync.Mutex
	initSig   chan os.Signal
	flushSig  chan *sync.WaitGroup
	closed    int32
	closeWait *sync.WaitGroup
	adapters  map[string]logAdapter
)

// logAdapter 定义了日志适配器接口。
// 实现此接口的类型可以作为日志输出的目标，如控制台、文件等。
type logAdapter interface {
	// init 初始化日志适配器。
	// 输入配置信息，返回此适配器支持的最大日志级别。
	init(prefs XPrefs.IBase) LevelType

	// write 写入一条日志记录。
	// 输入日志数据，返回写入过程中可能发生的错误。
	write(log *logData) error

	// flush 将缓冲区中的日志立即写入底层存储。
	flush()

	// close 关闭日志适配器，释放相关资源。
	close()
}

func init() { setup(XPrefs.Asset()) }

// setup 初始化日志系统。
// 输入配置信息，根据配置设置日志系统的运行参数。
// 此函数会初始化所有配置的日志适配器，并启动日志处理协程。
func setup(prefs XPrefs.IBase) {
	if prefs == nil {
		Panic("XLog.Init: prefs is nil.")
		return
	}

	initMu.Lock()
	defer initMu.Unlock()

	Close()
	atomic.SwapInt32(&closed, 0)
	closeWait = &sync.WaitGroup{}
	adapters = make(map[string]logAdapter)
	flushSig = make(chan *sync.WaitGroup, 1)

	levelMax = LevelUndefined
	for _, key := range prefs.Keys() {
		if !strings.HasPrefix(key, "Log/") {
			continue
		}
		name := strings.Split(key, "Log/")[1]
		if _, ok := adapters[name]; ok {
			Error("XLog.Init: dumplicated adapter: %v.", name)
			continue
		}

		var adapter logAdapter
		if name == "Std" {
			adapter = newStdAdapter()
		} else if name == "File" {
			adapter = newFileAdapter()
		} else {
			Warn("XLog.Init: unsupported adapter: %v.", name)
		}

		if adapter != nil {
			conf := prefs.Get(key).(XPrefs.IBase)
			level := adapter.init(conf)
			if level > levelMax {
				levelMax = level
			}
			adapters[name] = adapter
		}
	}

	initSig = make(chan os.Signal, 1)
	signal.Notify(initSig, syscall.SIGTERM, syscall.SIGINT)

	closeWait.Add(1)
	quit.GetWaiter().Add(1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()

		defer func() {
			for {
				if len(logCache) > 0 {
					log := <-logCache
					for name, adapter := range adapters {
						err := adapter.write(log)
						if err != nil {
							fmt.Fprintf(os.Stderr, "XLog.Listen: write in close: %v, error: %v\n", name, err)
						}
					}
					logPool.Put(log)
					continue
				} else {
					break
				}
			}
			for _, adapter := range adapters {
				adapter.flush()
				adapter.close()
			}
			closeWait.Done()
			quit.GetWaiter().Done()
		}()

		for {
			select {
			case log := <-logCache:
				for name, adapter := range adapters {
					err := adapter.write(log)
					if err != nil {
						fmt.Fprintf(os.Stderr, "XLog.Listen: write in queue: %v, error: %v\n", name, err)
					}
				}
				logPool.Put(log)
			case sig := <-flushSig:
				for len(logCache) > 0 {
					log := <-logCache
					for name, adapter := range adapters {
						err := adapter.write(log)
						if err != nil {
							fmt.Fprintf(os.Stderr, "XLog.Listen: write in flush: %v, error: %v\n", name, err)
						}
					}
					logPool.Put(log)
				}
				for _, adapter := range adapters {
					adapter.flush()
				}
				sig.Done()
			case sig, ok := <-initSig:
				if ok {
					fmt.Printf("XLog.Listen: receive signal of %v.\n", sig.String())
				} else {
					fmt.Printf("XLog.Listen: channel of signal is closed.\n")
				}
				return
			case <-quit.GetQuitChannel():
				fmt.Println("XLog.Listen: receive signal of QUIT.")
				return
			}
		}
	}()
	wg.Wait()
}

// Flush 将缓冲区中的所有日志立即写入到目标位置。
// 此函数会等待所有缓冲的日志被写入后才返回。
func Flush() {
	if initSig != nil && closed == 0 {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		select {
		case flushSig <- wg:
			wg.Wait()
		}
		Notice("XLog.Flush: logger has been flushed.")
	}
}

// Close 关闭日志系统并释放相关资源。
// 此函数会等待所有待写入的日志完成写入，然后关闭所有日志适配器。
func Close() {
	if initSig != nil && atomic.CompareAndSwapInt32(&closed, 0, 1) {
		signal.Stop(initSig)
		close(initSig)
		closeWait.Wait()
		Notice("XLog.Close: logger has been closed.")
	}
}

// Level 返回当前允许输出的最大日志级别。
// 任何高于此级别的日志都不会被记录，除非通过日志标签强制指定了更高的级别。
func Level() LevelType { return levelMax }

// Able 检查指定的日志级别是否允许输出。
// 输入日志级别，如果该级别的日志允许输出则返回 true，否则返回 false。
// 注意：如果存在日志标签且定义了级别，则标签的级别优先于全局级别。
func Able(level LevelType) bool {
	ret, _, _, _ := condition(level, nil)
	return ret
}

// condition 检查给定的日志级别是否可以根据配置的最大级别输出，并解析参数中的 LogTag。
// 输入日志级别和格式参数，返回是否允许输出、是否强制输出、日志标签和处理后的参数列表。
// 注意：标签中定义的日志级别优先于全局最大日志级别。
func condition(level LevelType, args []any) (bool, bool, *LogTag, []any) {
	var tag *LogTag
	var nargs []any

	// 快速路径：如果没有参数，直接使用上下文 tag
	if len(args) == 0 {
		goto CHECK_CONTEXT_TAG
	}

	// 检查第一个参数是否为 LogTag
	if arg0, ok := args[0].(*LogTag); ok && arg0 != nil {
		tag = arg0
		nargs = args[1:] // 只在确实有 LogTag 时才切片
		if tag.Level() != LevelUndefined {
			return level <= tag.Level(), true, tag, nargs
		}
		return level <= levelMax, false, tag, nargs
	}

	// 参数中没有 LogTag，保持原参数不变
	nargs = args

CHECK_CONTEXT_TAG:
	// 检查上下文 tag
	if ctxTag := Tag(); ctxTag != nil && ctxTag.Level() != LevelUndefined {
		return level <= ctxTag.Level(), true, ctxTag, nargs
	}
	return level <= levelMax, false, nil, nargs
}

// Panic 记录一条紧急日志，并触发 panic。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数会在记录日志后立即触发 panic，中断程序执行。
func Panic(data any, args ...any) {
	if data != nil {
		str := formatLog(data, args...)
		panic(str)
	}
}

// Emergency 记录一条紧急级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录导致系统完全不可用的灾难性故障。
func Emergency(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelEmergency, args); able {
		Print(LevelEmergency, force, tag, data, nargs...)
	}
}

// Alert 记录一条警报级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录需要立即引起注意和处理的系统状况。
func Alert(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelAlert, args); able {
		Print(LevelAlert, force, tag, data, nargs...)
	}
}

// Critical 记录一条严重级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录需要立即注意的严重系统故障。
func Critical(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelCritical, args); able {
		Print(LevelCritical, force, tag, data, nargs...)
	}
}

// Error 记录一条错误级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录需要解决的错误状况。
func Error(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelError, args); able {
		Print(LevelError, force, tag, data, nargs...)
	}
}

// Warn 记录一条警告级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录可能导致错误的潜在问题。
func Warn(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelWarn, args); able {
		Print(LevelWarn, force, tag, data, nargs...)
	}
}

// Notice 记录一条通知级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录值得注意但不一定是问题的事件。
func Notice(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelNotice, args); able {
		Print(LevelNotice, force, tag, data, nargs...)
	}
}

// Info 记录一条信息级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录系统的常规操作信息。
func Info(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelInfo, args); able {
		Print(LevelInfo, force, tag, data, nargs...)
	}
}

// Debug 记录一条调试级别的日志。
// 输入日志内容和可选的格式化参数，支持通过 LogTag 指定特定的日志级别和标签。
// 此函数用于记录系统调试和故障排除的详细信息。
func Debug(data any, args ...any) {
	if able, force, tag, nargs := condition(LevelDebug, args); able {
		Print(LevelDebug, force, tag, data, nargs...)
	}
}

// Print 记录一条指定级别的日志。
// 输入日志级别、是否强制输出、日志标签、日志内容和可选的格式化参数。
// 此函数是所有日志记录函数的底层实现，支持完整的日志记录功能。
func Print(level LevelType, force bool, tag *LogTag, data any, args ...any) {
	log := logPool.Get().(*logData)
	log.reset()
	log.level = level
	log.force = force
	log.data = data
	log.time = time.Now()
	if tag != nil {
		log.tag = tag.Text()
	}
	log.args = args

	if initSig == nil || closed == 1 {
		h, _, _ := formatTime(log.time)
		fmt.Println(string(append(h, log.text(true)...)))
	} else {
		logCache <- log
	}
}

// Size 返回当前日志缓冲区中的日志数量。
// 此函数可用于监控日志系统的积压情况。
func Size() int { return len(logCache) }

// logData 定义了一条日志记录的完整信息。
// 包含日志的级别、内容、标签、时间戳等元数据。
type logData struct {
	// level 存储日志的严重级别。
	level LevelType

	// force 标记是否强制写入日志，忽略级别限制。
	force bool

	// data 存储日志的主要内容。
	data any

	// args 存储用于格式化日志内容的参数列表。
	args []any

	// tag 存储日志的标签信息。
	tag string

	// time 记录日志产生的时间戳。
	time time.Time
}

// text 生成日志记录的文本表示。
// 输入是否包含标签信息，返回格式化后的日志文本。
// 如果启用了标签且存在标签信息，则在日志文本中包含标签。
func (log *logData) text(tag bool) string {
	if tag && log.tag != "" {
		return levelLabel[log.level] + " " + log.tag + " " + formatLog(log.data, log.args...)
	} else {
		return levelLabel[log.level] + " " + formatLog(log.data, log.args...)
	}
}

// reset 重置日志记录的所有字段为零值。
// 此方法在将日志对象放回对象池前调用，避免内存泄漏。
func (log *logData) reset() {
	log.level = LevelUndefined
	log.force = false
	log.data = nil
	log.args = nil
	log.tag = ""
}
