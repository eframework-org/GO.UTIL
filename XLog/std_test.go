// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/eframework-org/GO.UTIL/XPrefs"
)

// 测试 stdAdapter 的 init 方法。
func TestStdAdapterInit(t *testing.T) {
	// Test case 1: Custom configuration
	prefs := XPrefs.New()
	prefs.Set(stdPrefsLevel, LevelInfoStr)
	prefs.Set(stdPrefsColor, true)
	adapter := &stdAdapter{}
	level := adapter.init(prefs)
	if level != LevelInfo {
		t.Errorf("Expected Level to be LevelInfo, got %v", level)
	}
	if !adapter.color {
		t.Errorf("Expected color to be true, got %v", adapter.color)
	}

	// Test case 2: Invalid Level
	prefs = XPrefs.New()
	prefs.Set(stdPrefsLevel, "INVALID")
	adapter = &stdAdapter{}
	level = adapter.init(prefs)
	if level != LevelUndefined {
		t.Errorf("Expected Level to be LevelUndefined, got %v", level)
	}

	// Test case 3: Missing Level
	prefs = XPrefs.New()
	prefs.Set(stdPrefsColor, true)
	adapter = &stdAdapter{}
	level = adapter.init(prefs)
	if level != LevelInfo {
		t.Errorf("Expected Level to be LevelInfo, got %v", level)
	}
	if !adapter.color {
		t.Errorf("Expected color to be true, got %v", adapter.color)
	}
}

// 测试 stdAdapter 的 write 方法。
func TestStdAdapterWrite(t *testing.T) {
	// Create a buffer to capture the output
	buf := &bytes.Buffer{}

	// Create a stdAdapter instance
	adapter := &stdAdapter{
		writer: buf,
		level:  LevelInfo,
		color:  false,
	}

	// Test case 1: Log level is lower than adapter level and force is false
	log := &logData{
		level: LevelDebug,
		force: false,
		time:  time.Now(),
		data:  "Debug message",
	}
	err := adapter.write(log)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("Expected empty buffer, got %v", buf.String())
	}

	// Test case 2: Log level is higher than adapter level and force is false
	log.level = LevelError
	err = adapter.write(log)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("Expected non-empty buffer, got empty")
	}
	expected := formatStdAdapterLog(adapter, log)
	if buf.String() != expected {
		t.Errorf("Expected %v, got %v", expected, buf.String())
	}
	buf.Reset()

	// Test case 3: Log level is lower than adapter level but force is true
	log.level = LevelDebug
	log.force = true
	err = adapter.write(log)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("Expected non-empty buffer, got empty")
	}
	expected = formatStdAdapterLog(adapter, log)
	if buf.String() != expected {
		t.Errorf("Expected %v, got %v", expected, buf.String())
	}
	buf.Reset()

	// Test case 4: Color is enabled
	adapter.color = true
	err = adapter.write(log)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("Expected non-empty buffer, got empty")
	}
	expected = formatStdAdapterLog(adapter, log)
	if buf.String() != expected {
		t.Errorf("Expected %v, got %v", expected, buf.String())
	}
	buf.Reset()
}

// 测试 stdAdapter 的 flush 方法。
func TestStdAdapterFlush(t *testing.T) {
	adapter := &stdAdapter{}
	adapter.flush()
	// Since flush is a no-op, we just ensure it doesn't panic or cause any issues
}

// 测试 stdAdapter 的 close 方法。
func TestStdAdapterClose(t *testing.T) {
	adapter := &stdAdapter{}
	adapter.close()
	// Since close is a no-op, we just ensure it doesn't panic or cause any issues
}

func formatStdAdapterLog(apt *stdAdapter, log *logData) string {
	if log == nil {
		return ""
	}
	str := log.text(true)
	if apt.color {
		str = strings.Replace(str, levelLabel[log.level], stdBrushes[log.level](levelLabel[log.level]), 1)
	}
	h, _, _ := formatTime(log.time)
	return string(append(append(h, str...), '\n'))
}
