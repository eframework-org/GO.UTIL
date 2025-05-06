// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/eframework-org/GO.UTIL/XEnv"
	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/eframework-org/GO.UTIL/XString"
)

// 文件日志适配器的配置项及其默认值
const (
	prefsFileLevel          = "Level"                 // 日志输出级别
	prefsFileLevelDefault   = LevelNoticeStr          // 默认为 Notice 级别
	prefsFileRotate         = "Rotate"                // 是否启用日志文件轮转
	prefsFileRotateDefault  = true                    // 默认启用轮转
	prefsFileDaily          = "Daily"                 // 是否按天轮转
	prefsFileDailyDefault   = true                    // 默认按天轮转
	prefsFileMaxDay         = "MaxDay"                // 日志文件保留天数
	prefsFileMaxDayDefault  = 7                       // 默认保留7天
	prefsFileHourly         = "Hourly"                // 是否按小时轮转
	prefsFileHourlyDefault  = true                    // 默认按小时轮转
	prefsFileMaxHour        = "MaxHour"               // 日志文件保留小时数
	prefsFileMaxHourDefault = 168                     // 默认保留168小时（7天）
	prefsFilePath           = "Path"                  // 日志文件存储路径
	prefsFilePathDefault    = "${Env.LocalPath}/Log/" // 默认存储在本地Log目录
	prefsFileMaxFile        = "MaxFile"               // 最大文件数量
	prefsFileMaxFileDefault = 100                     // 默认保留100个文件
	prefsFileMaxLine        = "MaxLine"               // 单文件最大行数
	prefsFileMaxLineDefault = 1000000                 // 默认单文件100万行
	prefsFileMaxSize        = "MaxSize"               // 单文件最大体积（字节）
	prefsFileMaxSizeDefault = 1 << 27                 // 默认128MB
)

// fileAdapter 实现基于文件的日志输出适配器，支持按大小、行数、时间进行日志文件轮转。
// 可以配置日志级别、轮转策略、文件路径等参数，并自动清理过期的日志文件。
type fileAdapter struct {
	sync.RWMutex           // 保护并发访问的互斥锁
	level        LevelType // 日志输出级别
	rotate       bool      // 是否启用日志轮转
	daily        bool      // 是否按天轮转
	maxDay       int       // 日志保留天数
	hourly       bool      // 是否按小时轮转
	maxHour      int       // 日志保留小时数
	path         string    // 日志文件路径
	maxFile      int       // 最大文件数量
	maxLine      int       // 单文件最大行数
	maxSize      int       // 单文件最大字节数

	fileWriter     *os.File  // 当前日志文件的写入器
	curMaxLine     int       // 当前文件已写入的行数
	curMaxFile     int       // 当前日志文件数量
	curMaxSize     int       // 当前文件已写入的字节数
	dailyOpenDate  int       // 当前日志文件的创建日期
	dailyOpenTime  time.Time // 当前日志文件的创建时间（按天）
	hourlyOpenDate int       // 当前日志文件的创建小时
	hourlyOpenTime time.Time // 当前日志文件的创建时间（按小时）
	prefix, suffix string    // 日志文件名的前缀和后缀
}

// newFileAdapter 创建一个新的文件日志适配器实例。
// 返回的适配器需要通过 init 方法进行初始化后才能使用。
func newFileAdapter() *fileAdapter {
	apt := &fileAdapter{}
	return apt
}

// init 使用提供的配置初始化文件日志适配器。
// 设置日志级别、轮转策略、文件路径等参数，并创建必要的目录和文件。
// prefs 为配置参数，包含日志级别、轮转设置等。
// 返回配置的日志级别。
func (apt *fileAdapter) init(prefs XPrefs.IBase) LevelType {
	if prefs == nil {
		return LevelUndefined
	}
	tmpLevel := prefs.GetString(prefsFileLevel, prefsFileLevelDefault)
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
	apt.rotate = prefs.GetBool(prefsFileRotate, prefsFileRotateDefault)
	apt.daily = prefs.GetBool(prefsFileDaily, prefsFileDailyDefault)
	apt.maxDay = prefs.GetInt(prefsFileMaxDay, prefsFileMaxDayDefault)
	apt.hourly = prefs.GetBool(prefsFileHourly, prefsFileHourlyDefault)
	apt.maxHour = prefs.GetInt(prefsFileMaxHour, prefsFileMaxHourDefault)
	apt.path = filepath.Clean(XString.Eval(prefs.GetString(prefsFilePath, prefsFilePathDefault), XEnv.Vars(), XPrefs.Asset()))
	apt.maxFile = prefs.GetInt(prefsFileMaxFile, prefsFileMaxFileDefault)
	apt.maxLine = prefs.GetInt(prefsFileMaxLine, prefsFileMaxLineDefault)
	apt.maxSize = prefs.GetInt(prefsFileMaxSize, prefsFileMaxSizeDefault)

	// 处理路径逻辑
	if filepath.Ext(apt.path) == "" {
		// 如果路径没有扩展名，认为是目录
		apt.suffix = ".log"
		apt.prefix = ""
		// 确保路径以分隔符结尾
		apt.path = filepath.Join(apt.path, apt.suffix)
	} else {
		apt.suffix = filepath.Ext(apt.path)
		base := filepath.Base(apt.path)
		if base == apt.suffix {
			// 如果基本名称就是后缀（如 .log），则没有文件名
			apt.prefix = ""
		} else {
			// 正常的文件名情况
			apt.prefix = strings.TrimSuffix(base, apt.suffix)
			// 确保前缀不以点结尾
			apt.prefix = strings.TrimSuffix(apt.prefix, ".")
		}
	}

	apt.curMaxFile = 0

	err := apt.startLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fileAdapter.init(%q): %s\n", apt.path, err)
	}
	return apt.level
}

// flush 将缓冲区中的日志数据立即写入磁盘。
func (apt *fileAdapter) flush() {
	apt.fileWriter.Sync()
}

// close 关闭文件日志适配器。
// 将缓冲区中的数据写入文件并关闭文件句柄。
func (apt *fileAdapter) close() {
	apt.fileWriter.Close()
}

// startLogger 启动日志记录器。
// 创建或打开日志文件，初始化文件描述符。
// 如果启动失败，返回错误信息。
func (apt *fileAdapter) startLogger() error {
	file, err := apt.createLogFile()
	if err != nil {
		return err
	}
	if apt.fileWriter != nil {
		apt.fileWriter.Close()
	}
	apt.fileWriter = file
	return apt.initFd()
}

// needRotateDaily 检查是否需要按天轮转日志文件。
// 当达到最大行数、最大大小或跨天时返回 true。
// day 为当前日期。
func (apt *fileAdapter) needRotateDaily(day int) bool {
	return (apt.maxLine > 0 && apt.curMaxLine >= apt.maxLine) ||
		(apt.maxSize > 0 && apt.curMaxSize >= apt.maxSize) ||
		(apt.daily && day != apt.dailyOpenDate)
}

// needRotateHourly 检查是否需要按小时轮转日志文件。
// 当达到最大行数、最大大小或跨小时时返回 true。
// hour 为当前小时。
func (apt *fileAdapter) needRotateHourly(hour int) bool {
	return (apt.maxLine > 0 && apt.curMaxLine >= apt.maxLine) ||
		(apt.maxSize > 0 && apt.curMaxSize >= apt.maxSize) ||
		(apt.hourly && hour != apt.hourlyOpenDate)
}

// write 将日志数据写入文件。
// 在写入前检查是否需要轮转日志文件，支持按小时或按天轮转。
// log 为要写入的日志数据。
// 如果写入失败，返回错误信息。
func (apt *fileAdapter) write(log *logData) error {
	if log == nil {
		return errors.New("nil log")
	}
	if log.level > apt.level && !log.force {
		return nil
	}
	str := log.text(true)
	hd, d, h := formatTime(log.time)
	str = string(hd) + str + "\n"
	if apt.rotate {
		apt.RLock()
		if apt.needRotateHourly(h) {
			apt.RUnlock()
			apt.Lock()
			if apt.needRotateHourly(h) {
				if err := apt.doRotate(log.time); err != nil {
					fmt.Fprintf(os.Stderr, "fileAdapter.write(%q): %s\n", apt.path, err)
				}
			}
			apt.Unlock()
		} else if apt.needRotateDaily(d) {
			apt.RUnlock()
			apt.Lock()
			if apt.needRotateDaily(d) {
				if err := apt.doRotate(log.time); err != nil {
					fmt.Fprintf(os.Stderr, "fileAdapter.write(%q): %s\n", apt.path, err)
				}
			}
			apt.Unlock()
		} else {
			apt.RUnlock()
		}
	}

	apt.Lock() // 保证写入时序
	_, err := apt.fileWriter.Write([]byte(str))
	if err == nil {
		apt.curMaxLine++
		apt.curMaxSize += len(str)
	}
	apt.Unlock()
	return err
}

// createLogFile 创建或打开日志文件。
// 创建必要的目录结构，设置适当的文件权限。
// 返回文件句柄和可能的错误信息。
func (apt *fileAdapter) createLogFile() (*os.File, error) {
	// 使用 filepath 而不是 path 来处理跨平台路径
	dirPath := filepath.Dir(apt.path)

	// 创建所有必要的父目录，使用 0755 权限
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directories: %v", err)
	}

	// 打开或创建文件
	fd, err := os.OpenFile(apt.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	// 设置文件权限
	if err = os.Chmod(apt.path, 0644); err != nil {
		fd.Close()
		return nil, fmt.Errorf("failed to set file permissions: %v", err)
	}

	return fd, nil
}

// initFd 初始化文件描述符。
// 获取文件信息，设置文件大小和时间相关的计数器。
// 如果启用了轮转，启动相应的轮转协程。
func (apt *fileAdapter) initFd() error {
	fd := apt.fileWriter
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s", err)
	}
	apt.curMaxSize = int(fInfo.Size())
	apt.dailyOpenTime = time.Now()
	apt.dailyOpenDate = apt.dailyOpenTime.Day()
	apt.hourlyOpenTime = time.Now()
	apt.hourlyOpenDate = apt.hourlyOpenTime.Hour()
	apt.curMaxLine = 0
	if apt.hourly {
		go apt.hourlyRotate(apt.hourlyOpenTime)
	} else if apt.daily {
		go apt.dailyRotate(apt.dailyOpenTime)
	}
	if fInfo.Size() > 0 && apt.maxLine > 0 {
		count, err := apt.lines()
		if err != nil {
			return err
		}
		apt.curMaxLine = count
	}
	return nil
}

// dailyRotate 执行按天轮转的定时任务。
// 计算下一天的时间点，定时检查并执行轮转。
// openTime 为当前日志文件的打开时间。
func (apt *fileAdapter) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	apt.Lock()
	if apt.needRotateDaily(time.Now().Day()) {
		if err := apt.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", apt.path, err)
		}
	}
	apt.Unlock()
}

// hourlyRotate 执行按小时轮转的定时任务。
// 计算下一小时的时间点，定时检查并执行轮转。
// openTime 为当前日志文件的打开时间。
func (apt *fileAdapter) hourlyRotate(openTime time.Time) {
	y, m, d := openTime.Add(1 * time.Hour).Date()
	h, _, _ := openTime.Add(1 * time.Hour).Clock()
	nextHour := time.Date(y, m, d, h, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextHour.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	apt.Lock()
	if apt.needRotateHourly(time.Now().Hour()) {
		if err := apt.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", apt.path, err)
		}
	}
	apt.Unlock()
}

// lines 计算日志文件的总行数。
// 通过读取文件内容并计数换行符来统计行数。
// 返回行数和可能的错误信息。
func (apt *fileAdapter) lines() (int, error) {
	fd, err := os.Open(apt.path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

// doRotate 执行日志文件轮转。
// 关闭当前文件，创建新文件，并更新相关计数器。
// logTime 为触发轮转的日志时间。
// 如果轮转过程中发生错误，返回错误信息。
func (apt *fileAdapter) doRotate(logTime time.Time) error {
	// file exists
	// Find the next available number
	num := apt.curMaxFile + 1
	fName := ""
	format := ""
	var openTime time.Time

	_, err := os.Lstat(apt.path)
	if err != nil {
		// even if the file is not exist or other, we should RESTART the logger
		goto RESTART_LOGGER
	}

	if apt.hourly {
		format = "2006-01-02-15"
		openTime = apt.hourlyOpenTime
	} else if apt.daily {
		format = "2006-01-02"
		openTime = apt.dailyOpenTime
	}

	// 生成轮转文件名
	if apt.maxLine > 0 || apt.maxSize > 0 {
		for ; err == nil && num <= apt.maxFile; num++ {
			if apt.prefix == "" {
				// 无文件名情况：使用序号作为文件名
				fName = filepath.Join(filepath.Dir(apt.path), fmt.Sprintf("%03d%s",
					num,
					apt.suffix))
			} else {
				// 有文件名情况：在原文件名后添加序号
				if format != "" {
					// 按时间轮转
					fName = filepath.Join(filepath.Dir(apt.path),
						fmt.Sprintf("%s.%s.%03d%s",
							strings.TrimSuffix(apt.prefix, "."),
							logTime.Format(format),
							num,
							apt.suffix))
				} else {
					// 按行数或大小轮转
					fName = filepath.Join(filepath.Dir(apt.path),
						fmt.Sprintf("%s.%03d%s",
							strings.TrimSuffix(apt.prefix, "."),
							num,
							apt.suffix))
				}
			}
			_, err = os.Lstat(fName)
		}
	} else {
		if apt.prefix == "" {
			fName = filepath.Join(filepath.Dir(apt.path), fmt.Sprintf("%s.%03d%s",
				openTime.Format(format),
				num,
				apt.suffix))
		} else {
			fName = filepath.Join(filepath.Dir(apt.path),
				fmt.Sprintf("%s.%s.%03d%s",
					strings.TrimSuffix(apt.prefix, "."),
					openTime.Format(format),
					num,
					apt.suffix))
		}
		_, err = os.Lstat(fName)
		apt.curMaxFile = num
	}

	// return error if the last file checked still existed
	if err == nil {
		return fmt.Errorf("rotate error: cannot find free log number to rename %s", apt.path)
	}

	// close fileWriter before rename
	apt.fileWriter.Close()

	// Rename the file to its new found name
	// even if occurs error, we MUST guarantee to restart new logger
	err = os.Rename(apt.path, fName)
	if err != nil {
		goto RESTART_LOGGER
	}

RESTART_LOGGER:
	startLoggerErr := apt.startLogger()
	go apt.deleteOld()

	if startLoggerErr != nil {
		return fmt.Errorf("rotate start error: %s", startLoggerErr)
	}
	if err != nil {
		return fmt.Errorf("rotate error: %s", err)
	}
	return nil
}

// deleteOld 清理过期的日志文件。
// 遍历日志目录，根据文件修改时间和配置的保留策略删除过期文件。
func (apt *fileAdapter) deleteOld() {
	dir := filepath.Dir(apt.path)
	var basePattern string
	if apt.prefix == "" {
		// 无文件名情况：匹配所有时间戳格式的日志文件
		basePattern = fmt.Sprintf("*[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]*%s", apt.suffix)
	} else {
		// 有文件名情况：匹配指定文件名的日志文件
		basePattern = fmt.Sprintf("%s.*%s", apt.prefix, apt.suffix)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil || err != nil {
			return
		}

		matched, err := filepath.Match(basePattern, filepath.Base(path))
		if err != nil || !matched {
			return
		}

		if apt.hourly {
			if !info.IsDir() && info.ModTime().Add(1*time.Hour*time.Duration(apt.maxHour)).Before(time.Now()) {
				os.Remove(path)
			}
		} else if apt.daily {
			if !info.IsDir() && info.ModTime().Add(24*time.Hour*time.Duration(apt.maxDay)).Before(time.Now()) {
				os.Remove(path)
			}
		}
		return
	})
}
