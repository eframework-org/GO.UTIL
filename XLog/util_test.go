// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Test formatTime function
func TestFormatTime(t *testing.T) {
	now := time.Now()
	h, d, _ := formatTime(now)
	expected := fmt.Sprintf("[%02d/%02d %02d:%02d:%02d.%03d] ", now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000)
	if string(h) != expected {
		t.Errorf("Expected %v, got %v", expected, string(h))
	}
	if d != int(now.Day()) {
		t.Errorf("Expected day %v, got %v", now.Day(), d)
	}
}

// Test formatLog function
func TestFormatLog(t *testing.T) {
	str := formatLog("Hello %s", "World")
	if str != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", str)
	}

	str = formatLog("Hello World")
	if str != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", str)
	}

	str = formatLog(123, "Hello", "World")
	if str != "123 Hello World" {
		t.Errorf("Expected '123 Hello World', got %v", str)
	}
}

// Test UnAddr function
func TestUnAddr(t *testing.T) {
	input := "192.168.1.1:8080 example.com:443"
	expected := "**.**.**.**:8080 ***.***:443"
	result := UnAddr(input)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Test Caller function
func TestCaller(t *testing.T) {
	caller := Caller(0, false)
	if !strings.Contains(caller, "TestCaller") {
		t.Errorf("Expected caller to contain 'TestCaller', got %v", caller)
	}

	caller = Caller(0, true)
	if !strings.Contains(caller, "util_test.go") {
		t.Errorf("Expected caller to contain 'util_test.go', got %v", caller)
	}
}

// Test Trace function
func TestTrace(t *testing.T) {
	trace, count := Trace(0, errors.New("Test error"))
	if !strings.Contains(trace, "Test error") {
		t.Errorf("Expected trace to contain 'Test error', got %v", trace)
	}
	if count < 1 {
		t.Errorf("Expected count to be at least 1, got %v", count)
	}
}

// Test Elapse function
func TestElapse(t *testing.T) {
	start := time.Now()
	defer Elapse(0)()
	end := time.Now()
	if end.Before(start) {
		t.Errorf("Expected end time to be after start time")
	}
}

func TestCaught(t *testing.T) {
	handlerCalled := false
	handler := func(s string, i int) {
		handlerCalled = true
	}
	func() {
		defer Caught(false, handler)
		panic("test panic")
	}()
	if !handlerCalled {
		t.Errorf("Expected handler to be called")
	}
}
