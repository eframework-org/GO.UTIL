// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
)

// 测试日志级别常量和 Able 函数.
func TestLogLevels(t *testing.T) {
	tests := []struct {
		level    LevelType
		expected bool
	}{
		{LevelEmergency, true},
		{LevelAlert, true},
		{LevelCritical, true},
		{LevelError, true},
		{LevelWarn, true},
		{LevelNotice, true},
		{LevelInfo, true},
		{LevelDebug, true},
		{LevelType(100), false}, // Invalid level
	}

	// Set maximum level to Debug to allow all valid levels
	levelMax = LevelDebug

	for _, test := range tests {
		result := Able(test.level)
		if result != test.expected {
			t.Errorf("Able(%v) = %v; want %v", test.level, result, test.expected)
		}
	}
}

// 测试日志系统的初始化和关闭功能.
func TestInitAndClose(t *testing.T) {
	// Create temporary file for testing
	tmpFile, err := os.CreateTemp("", "log_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create configuration
	prefs := XPrefs.New()
	fileConf := XPrefs.New()
	fileConf.Set(prefsFilePath, tmpFile.Name())
	fileConf.Set(prefsFileLevel, LevelDebugStr)
	fileConf.Set(prefsFileRotate, false) // 禁用日志轮转
	prefs.Set("Log/File", fileConf)

	// Test initialization
	setup(prefs)
	if len(adapters) != 1 {
		t.Errorf("Expected 1 adapter, got %d", len(adapters))
	}

	// Test logging
	message := "Test message"
	Info(message)

	// 确保日志写入完成
	Flush()
	time.Sleep(100 * time.Millisecond)

	// 读取日志内容
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 验证日志内容
	if !bytes.Contains(content, []byte(message)) {
		t.Errorf("Log file does not contain expected message: %s\nActual content: %s", message, string(content))
	}

	// 关闭日志系统
	Close()
}

// 测试所有日志记录方法.
func TestLogMethods(t *testing.T) {
	// Create temporary file for testing
	tmpFile, err := os.CreateTemp("", "log_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create configuration
	prefs := XPrefs.New()
	fileConf := XPrefs.New()
	fileConf.Set(prefsFilePath, tmpFile.Name())
	fileConf.Set(prefsFileLevel, LevelDebugStr)
	fileConf.Set(prefsFileRotate, false) // 禁用日志轮转
	prefs.Set("Log/File", fileConf)

	// Initialize logger
	setup(prefs)

	tests := []struct {
		level   LevelType
		logFunc func(data any, args ...any)
		message string
	}{
		{LevelEmergency, Emergency, "Emergency message"},
		{LevelAlert, Alert, "Alert message"},
		{LevelCritical, Critical, "Critical message"},
		{LevelError, Error, "Error message"},
		{LevelWarn, Warn, "Warning message"},
		{LevelNotice, Notice, "Notice message"},
		{LevelInfo, Info, "Info message"},
		{LevelDebug, Debug, "Debug message"},
	}

	// 写入所有日志消息
	for _, test := range tests {
		test.logFunc(test.message)
	}

	// 确保日志写入完成
	Flush()
	time.Sleep(100 * time.Millisecond)

	// 验证日志内容
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 检查每条消息是否都被正确记录
	for _, test := range tests {
		if !bytes.Contains(content, []byte(test.message)) {
			t.Errorf("Log file does not contain expected message: %s\nActual content: %s", test.message, string(content))
		}
	}

	// 关闭日志系统
	Close()
}

// 测试 Panic 函数功能.
func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, got none")
		}
	}()

	Panic("Test panic")
}

// 测试并发日志记录操作.
func TestConcurrentLogging(t *testing.T) {
	// Create temporary file for testing
	tmpFile, err := os.CreateTemp("", "concurrent_log_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Initialize logger
	prefs := XPrefs.New()
	fileConf := XPrefs.New()
	fileConf.Set("Path", tmpFile.Name())
	fileConf.Set("Level", LevelDebug)
	fileConf.Set("Rotate", false) // 禁用日志轮转
	prefs.Set("Log/File", fileConf)

	setup(prefs)

	// Test concurrent logging
	var wg sync.WaitGroup
	numGoroutines := 10
	numLogs := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numLogs; j++ {
				Info(fmt.Sprintf("Goroutine %d: Log %d", id, j))
			}
		}(i)
	}
	wg.Wait()

	// Test flush and close
	Flush()
	Close()

	// Verify log file exists and has content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Error("Log file is empty")
	}
}

// 测试日志数据池功能.
func TestLogDataPool(t *testing.T) {
	// Get log data from pool
	log1 := logPool.Get().(*logData)
	if log1 == nil {
		t.Fatal("Failed to get log data from pool")
	}

	// Set some data
	log1.level = LevelInfo
	log1.data = "Test message"

	// Reset and return to pool
	log1.reset()
	logPool.Put(log1)

	// Get another log data and verify it's reset
	log2 := logPool.Get().(*logData)
	if log2.level != LevelUndefined || log2.data != nil {
		t.Error("Log data not properly reset")
	}
}

// 测试日志标签功能，包括级别控制和参数解析.
func TestLogWithTag(t *testing.T) {
	// 创建临时文件用于测试
	tmpFile, err := os.CreateTemp("", "logtag_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// 初始化日志系统
	prefs := XPrefs.New()
	fileConf := XPrefs.New()
	fileConf.Set(prefsFilePath, tmpFile.Name())
	fileConf.Set(prefsFileLevel, LevelDebugStr)
	fileConf.Set(prefsFileRotate, false)
	prefs.Set("Log/File", fileConf)
	setup(prefs)

	tests := []struct {
		name      string
		setup     func()
		logFunc   func()
		cleanup   func()
		expected  string
		shouldLog bool
	}{
		{
			name: "Tag in args overrides global level",
			setup: func() {
				levelMax = LevelInfo // 设置全局级别为 Info
			},
			logFunc: func() {
				tag := GetTag()
				tag.Level(LevelDebug)
				tag.Set("source", "test1")
				Debug("Debug message", tag)
			},
			cleanup:   func() {},
			expected:  "[source=test1] Debug message",
			shouldLog: true,
		},
		{
			name: "Context tag controls output",
			setup: func() {
				levelMax = LevelInfo
				tag := GetTag()
				tag.Level(LevelDebug)
				tag.Set("context", "test2")
				Watch(tag)
			},
			logFunc: func() {
				Debug("Context debug message")
			},
			cleanup: func() {
				Defer()
			},
			expected:  "[context=test2] Context debug message",
			shouldLog: true,
		},
		{
			name:  "Multiple tag values",
			setup: func() {},
			logFunc: func() {
				tag := GetTag()
				tag.Set("key1", "value1")
				tag.Set("key2", "value2")
				Info("Multi tag message", tag)
			},
			cleanup:   func() {},
			expected:  "[key1=value1, key2=value2] Multi tag message",
			shouldLog: true,
		},
		{
			name:  "Format string with tag",
			setup: func() {},
			logFunc: func() {
				tag := GetTag()
				tag.Set("format", "test")
				Info("Count: %d", tag, 42)
			},
			cleanup:   func() {},
			expected:  "[format=test] Count: 42",
			shouldLog: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()

			// 记录当前日志大小
			initialSize, err := os.Stat(tmpFile.Name())
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			initialBytes := int64(0)
			if initialSize != nil {
				initialBytes = initialSize.Size()
			}

			// 执行日志记录
			test.logFunc()

			// 确保日志写入完成
			Flush()
			time.Sleep(100 * time.Millisecond)

			// 检查文件大小是否变化
			currentSize, err := os.Stat(tmpFile.Name())
			if err != nil {
				t.Fatal(err)
			}
			hasNewContent := currentSize.Size() > initialBytes

			if hasNewContent != test.shouldLog {
				t.Errorf("Log output mismatch: expected shouldLog=%v, got hasNewContent=%v",
					test.shouldLog, hasNewContent)
			}

			if test.shouldLog {
				// 读取日志内容
				content, err := os.ReadFile(tmpFile.Name())
				if err != nil {
					t.Fatal(err)
				}

				// 验证日志内容
				if !bytes.Contains(content, []byte(test.expected)) {
					t.Errorf("Log content does not contain expected message.\nExpected: %s\nActual: %s",
						test.expected, string(content))
				}
			}

			test.cleanup()
		})
	}

	Close()
}
