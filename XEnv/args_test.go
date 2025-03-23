// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"os"
	"sync"
	"testing"
)

func resetCache() {
	argsCacheLock.Lock()
	argsCache = nil
	argsCacheInit = sync.Once{}
	argsCacheLock.Unlock()
}

func TestGetArg(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		resetCache()
	}()

	tests := []struct {
		name     string
		args     []string
		key      string
		expected string
	}{
		{
			name:     "equals format",
			args:     []string{"prog", "--test=value"},
			key:      "test",
			expected: "value",
		},
		{
			name:     "space format",
			args:     []string{"prog", "--test", "value"},
			key:      "test",
			expected: "value",
		},
		{
			name:     "key not found",
			args:     []string{"prog", "--other=value"},
			key:      "test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			resetCache() // Reset cache for each test case
			if got := GetArg(tt.key); got != tt.expected {
				t.Errorf("GetArg(%q) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestGetArgs(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		resetCache()
	}()

	tests := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			name: "mixed format",
			args: []string{"prog", "--key1=value1", "--key2", "value2"},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "empty args",
			args:     []string{"prog"},
			expected: map[string]string{},
		},
		{
			name: "invalid format ignored",
			args: []string{"prog", "--key1=value1", "invalid", "--key2=value2"},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			resetCache() // Reset cache for each test case
			got := GetArgs()

			if len(got) != len(tt.expected) {
				t.Errorf("GetArgs() got %v items, want %v items", len(got), len(tt.expected))
			}

			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("GetArgs()[%q] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		resetCache()
	}()

	os.Args = []string{"prog", "--test=value", "--concurrent=true"}

	var wg sync.WaitGroup
	numGoroutines := 10

	// Test concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			args := GetArgs()
			if args["test"] != "value" {
				t.Errorf("Concurrent GetArgs() failed")
			}
		}()
	}

	wg.Wait()
}

func TestCacheConsistency(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		resetCache()
	}()

	os.Args = []string{"prog", "--cache=test"}

	// First call should initialize cache
	args1 := GetArgs()

	// Modify os.Args
	os.Args = []string{"prog", "--different=value"}

	// Second call should return cached result
	args2 := GetArgs()

	if args1["cache"] != args2["cache"] {
		t.Error("Cache inconsistency detected")
	}
}
