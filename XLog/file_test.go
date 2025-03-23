// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
)

// 测试用例的通用结构
type testCase struct {
	name     string
	path     string
	expected struct {
		fileNameOnly string
		suffix       string
	}
}

// 测试使用配置初始化 fileAdapter
func TestFileAdapterInit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cases := []testCase{
		{
			name: "Normal Path",
			path: filepath.Join(tempDir, "test.log"),
			expected: struct {
				fileNameOnly string
				suffix       string
			}{
				fileNameOnly: "test",
				suffix:       ".log",
			},
		},
		{
			name: "With Suffix Only",
			path: filepath.Join(tempDir, ".log"),
			expected: struct {
				fileNameOnly string
				suffix       string
			}{
				fileNameOnly: "",
				suffix:       ".log",
			},
		},
		{
			name: "Directory Only",
			path: tempDir,
			expected: struct {
				fileNameOnly string
				suffix       string
			}{
				fileNameOnly: "",
				suffix:       ".log",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefs := XPrefs.New()
			prefs.Set(prefsFileLevel, LevelDebugStr)
			prefs.Set(prefsFilePath, tc.path)

			adapter := newFileAdapter()
			adapter.init(prefs)

			if adapter.prefix != tc.expected.fileNameOnly {
				t.Errorf("Expected fileNameOnly to be %q, got %q", tc.expected.fileNameOnly, adapter.prefix)
			}
			if adapter.suffix != tc.expected.suffix {
				t.Errorf("Expected suffix to be %q, got %q", tc.expected.suffix, adapter.suffix)
			}
		})
	}
}

// 测试日志写入和轮转
func TestFileAdapterWriteAndRotate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cases := []struct {
		name         string
		path         string
		rotate       bool
		hourly       bool
		daily        bool
		maxLine      int
		maxFile      int
		writeNum     int
		checkRotated bool
	}{
		{
			name:         "Normal Write Without Rotation",
			path:         filepath.Join(tempDir, "test1.log"),
			rotate:       false,
			writeNum:     3,
			checkRotated: false,
		},
		{
			name:         "Write With Line Rotation",
			path:         filepath.Join(tempDir, "test2.log"),
			rotate:       true,
			maxLine:      2,
			maxFile:      2,
			writeNum:     5,
			checkRotated: true,
		},
		{
			name:         "Write With Suffix Only",
			path:         filepath.Join(tempDir, ".customext"),
			rotate:       true,
			maxLine:      2,
			maxFile:      2,
			writeNum:     5,
			checkRotated: true,
		},
		{
			name:         "Write With Directory Only",
			path:         tempDir,
			rotate:       true,
			maxLine:      2,
			maxFile:      2,
			writeNum:     5,
			checkRotated: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefs := XPrefs.New()
			prefs.Set(prefsFileLevel, LevelDebugStr)
			prefs.Set(prefsFilePath, tc.path)
			prefs.Set(prefsFileRotate, tc.rotate)
			prefs.Set(prefsFileMaxLine, tc.maxLine)
			prefs.Set(prefsFileMaxFile, tc.maxFile)
			prefs.Set(prefsFileHourly, tc.hourly)
			prefs.Set(prefsFileDaily, tc.daily)

			adapter := newFileAdapter()
			adapter.init(prefs)

			// 写入日志
			logTime := time.Now()
			for i := 0; i < tc.writeNum; i++ {
				logData := &logData{
					level: LevelInfo,
					force: false,
					time:  logTime,
					data:  fmt.Sprintf("Test log message %d", i),
				}
				err := adapter.write(logData)
				if err != nil {
					t.Errorf("Write failed: %v", err)
				}
			}

			// 检查原始日志文件是否存在
			_, err := os.Stat(adapter.path)
			if err != nil {
				t.Errorf("Expected original log file to exist: %v", err)
			}

			if tc.checkRotated {
				// 检查轮转的日志文件
				dir := filepath.Dir(adapter.path)
				files, err := os.ReadDir(dir)
				if err != nil {
					t.Fatalf("Failed to read directory: %v", err)
				}

				rotatedCount := 0
				for _, file := range files {
					if !file.IsDir() {
						matched, err := filepath.Match(fmt.Sprintf("%s*%s", adapter.prefix, adapter.suffix), file.Name())
						if err == nil && matched {
							rotatedCount++
						}
					}
				}

				expectedFiles := tc.maxFile + 1 // 包括当前文件
				if rotatedCount < expectedFiles {
					t.Errorf("Expected at least %d log files, got %d", expectedFiles, rotatedCount)
				}
			}
		})
	}
}

// 测试旧日志文件的删除
func TestFileAdapterCleanup(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cases := []struct {
		name      string
		path      string
		hourly    bool
		maxHour   int
		daily     bool
		maxDay    int
		createOld bool
	}{
		{
			name:      "Hourly Cleanup",
			path:      filepath.Join(tempDir, "test1.log"),
			hourly:    true,
			maxHour:   1,
			createOld: true,
		},
		{
			name:      "Daily Cleanup",
			path:      filepath.Join(tempDir, "test2.log"),
			daily:     true,
			maxDay:    1,
			createOld: true,
		},
		{
			name:      "Suffix Only Cleanup",
			path:      filepath.Join(tempDir, ".log"),
			hourly:    true,
			maxHour:   1,
			createOld: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefs := XPrefs.New()
			prefs.Set(prefsFileLevel, LevelDebugStr)
			prefs.Set(prefsFilePath, tc.path)
			prefs.Set(prefsFileRotate, true)
			prefs.Set(prefsFileHourly, tc.hourly)
			prefs.Set(prefsFileMaxHour, tc.maxHour)
			prefs.Set(prefsFileDaily, tc.daily)
			prefs.Set(prefsFileMaxDay, tc.maxDay)

			adapter := newFileAdapter()
			adapter.init(prefs)

			if tc.createOld {
				// 创建旧的日志文件
				now := time.Now()
				oldTime := now.Add(-2 * time.Hour)
				var oldTimeStr string
				if tc.hourly {
					oldTimeStr = oldTime.Format("2006-01-02-15")
				} else if tc.daily {
					oldTime = now.AddDate(0, 0, -2) // 2天前
					oldTimeStr = oldTime.Format("2006-01-02")
				}

				var oldFiles []string
				if adapter.prefix == "" {
					oldFiles = []string{
						filepath.Join(filepath.Dir(tc.path), fmt.Sprintf("%s.001%s", oldTimeStr, adapter.suffix)),
						filepath.Join(filepath.Dir(tc.path), fmt.Sprintf("%s.002%s", oldTimeStr, adapter.suffix)),
					}
				} else {
					oldFiles = []string{
						filepath.Join(filepath.Dir(tc.path), fmt.Sprintf("%s.%s.001%s", adapter.prefix, oldTimeStr, adapter.suffix)),
						filepath.Join(filepath.Dir(tc.path), fmt.Sprintf("%s.%s.002%s", adapter.prefix, oldTimeStr, adapter.suffix)),
					}
				}

				// 创建文件并设置修改时间
				for _, file := range oldFiles {
					err := os.WriteFile(file, []byte("old log"), 0644)
					if err != nil {
						t.Fatalf("Failed to create file %s: %v", file, err)
					}
					err = os.Chtimes(file, oldTime, oldTime)
					if err != nil {
						t.Fatalf("Failed to set modification time for %s: %v", file, err)
					}
				}

				// 执行清理
				adapter.deleteOld()

				// 验证文件是否被删除
				for _, file := range oldFiles {
					_, err := os.Stat(file)
					if err == nil {
						t.Errorf("Expected %s to be deleted", file)
					}
				}
			}
		})
	}
}
