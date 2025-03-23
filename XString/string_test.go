// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XString

import (
	"encoding/base64"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestToInt 测试 ToInt 函数
func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Valid Integer", "123", 123},
		{"Negative Integer", "-456", -456},
		{"Zero", "0", 0},
		{"Invalid Input", "abc", 0},
		{"Empty String", "", 0},
		{"Mixed Content", "123abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToString 测试 ToString 函数
func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Positive Integer", 123, "123"},
		{"Negative Integer", -456, "-456"},
		{"Zero", 0, "0"},
		{"Large Number", 999999, "999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToFixed 测试 ToFixed 函数
func TestToFixed(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		decimals []int
		expected string
	}{
		{"Float32", float32(3.14159), []int{2}, "3.14"},
		{"Float64", 3.14159, []int{3}, "3.142"},
		{"Zero Decimals", 3.14159, []int{0}, "3"},
		{"Negative Decimals", 3.14159, []int{-2}, "3"},
		{"No Decimals Specified", 3.14159, nil, "3.14"},
		{"Invalid Type", "3.14", nil, ""},
		{"Nil Input", nil, nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToFixed(tt.input, tt.decimals...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSplit 测试 Split 函数
func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		sep      string
		expected []string
	}{
		{"Basic Split", "a,b,c", ",", []string{"a", "b", "c"}},
		{"Empty Parts", "a,,c", ",", []string{"a", "", "c"}},
		{"No Separator", "abc", ",", []string{"abc"}},
		{"Empty String", "", ",", []string{""}},
		{"Empty Separator", "abc", "", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Split(tt.str, tt.sep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIndexOf 测试 IndexOf 函数
func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected int
	}{
		{"Found Substring", "hello world", "world", 6},
		{"Substring at Start", "hello world", "hello", 0},
		{"Substring Not Found", "hello world", "xyz", -1},
		{"Empty Substring", "hello world", "", 0},
		{"Empty String", "", "world", -1},
		{"Both Empty", "", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IndexOf(tt.str, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSub 测试 Sub 函数
func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		from     int
		to       int
		expected string
	}{
		{"Valid Range", "hello world", 0, 5, "hello"},
		{"Mid Range", "hello world", 6, 11, "world"},
		{"Invalid Range", "hello world", 5, 3, ""},
		{"Negative Indices", "hello world", -1, 5, ""},
		{"Out of Bounds", "hello world", 0, 20, "hello world"},
		{"Empty String", "", 0, 5, ""},
		{"Unicode String", "你好世界", 0, 2, "你好"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sub(tt.str, tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToBuffer 测试 ToBuffer 函数
func TestToBuffer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{"ASCII String", "hello", []byte("hello")},
		{"Empty String", "", []byte("")},
		{"Unicode String", "你好", []byte("你好")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBuffer(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFromBuffer 测试 FromBuffer 函数
func TestFromBuffer(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"ASCII Bytes", []byte("hello"), "hello"},
		{"Empty Bytes", []byte(""), ""},
		{"Unicode Bytes", []byte("你好"), "你好"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromBuffer(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestEvaluator implements IEval for testing
type TestEvaluator struct {
	replacements map[string]string
}

func (e *TestEvaluator) Eval(input string) string {
	result := input
	for key, value := range e.replacements {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}

// TestEval 测试 Eval 函数
func TestEval(t *testing.T) {
	t.Run("Empty Input", func(t *testing.T) {
		result := Eval("")
		if result != "" {
			t.Errorf("Expected empty string, got %v", result)
		}
	})

	t.Run("With IEval", func(t *testing.T) {
		evaluator := &TestEvaluator{
			replacements: map[string]string{
				"${name}": "John",
				"${age}":  "30",
			},
		}

		input := "Name: ${name}, Age: ${age}"
		expected := "Name: John, Age: 30"

		result := Eval(input, evaluator)
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("With Map", func(t *testing.T) {
		replacements := map[string]string{
			"${city}":    "New York",
			"${country}": "USA",
		}

		input := "City: ${city}, Country: ${country}"
		expected := "City: New York, Country: USA"

		result := Eval(input, replacements)
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Multiple Sources", func(t *testing.T) {
		evaluator := &TestEvaluator{
			replacements: map[string]string{
				"${name}": "John",
			},
		}

		replacements := map[string]string{
			"${age}": "30",
		}

		input := "Name: ${name}, Age: ${age}"
		expected := "Name: John, Age: 30"

		result := Eval(input, evaluator, replacements)
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Nil Source", func(t *testing.T) {
		input := "Test string"
		result := Eval(input, nil)
		if result != input {
			t.Errorf("Expected %v, got %v", input, result)
		}
	})
}

// TestRandom 测试随机字符串生成功能
func TestRandom(t *testing.T) {
	t.Run("Format", func(t *testing.T) {
		tests := []struct {
			name       string
			format     string
			wantLength int
			wantFormat string
		}{
			{
				name:       "Default Format (N)",
				format:     "",
				wantLength: 32,
				wantFormat: "^[0-9a-f]{32}$",
			},
			{
				name:       "D Format",
				format:     "D",
				wantLength: 36,
				wantFormat: "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
			},
			{
				name:       "N Format",
				format:     "N",
				wantLength: 32,
				wantFormat: "^[0-9a-f]{32}$",
			},
			{
				name:       "Invalid Format",
				format:     "X",
				wantLength: 32,
				wantFormat: "^[0-9a-f]{32}$",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var result string
				if tt.format == "" {
					result = Random()
				} else {
					result = Random(tt.format)
				}

				// 检查长度
				assert.Equal(t, tt.wantLength, len(result), "unexpected length")

				// 检查格式
				matched, err := regexp.MatchString(tt.wantFormat, result)
				assert.NoError(t, err, "regexp.MatchString error")
				assert.True(t, matched, "Random format mismatch")

				// 检查版本号（应该是版本4）
				if tt.format == "D" && len(result) >= 14 {
					assert.Equal(t, "4", string(result[14]), "incorrect version")
				} else if tt.format != "D" && len(result) >= 12 {
					assert.Equal(t, "4", string(result[12]), "incorrect version")
				}

				// 检查变体位（如果是D格式）
				if tt.format == "D" {
					parts := strings.Split(result, "-")
					assert.Equal(t, 5, len(parts), "wrong UUID format")
					variant := parts[3][0]
					assert.Contains(t, "89ab", string(variant), "invalid variant")
				}
			})
		}
	})

	t.Run("Uniqueness", func(t *testing.T) {
		const count = 1000
		results := make(map[string]bool, count)

		for i := 0; i < count; i++ {
			result := Random()
			assert.False(t, results[result], "generated duplicate value: %v", result)
			results[result] = true
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		const goroutines = 100
		done := make(chan string, goroutines)
		results := make(map[string]bool)
		var mutex sync.Mutex
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				result := Random()
				done <- result
			}()
		}

		go func() {
			wg.Wait()
			close(done)
		}()

		for result := range done {
			mutex.Lock()
			assert.False(t, results[result], "generated duplicate value in concurrent execution: %v", result)
			results[result] = true
			mutex.Unlock()
		}
	})
}

// TestCrypt 测试加密和解密功能
func TestCrypt(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		tests := []struct {
			name     string
			data     any // string或[]byte
			key      string
			wantFail bool
		}{
			// 基本功能测试
			{
				name: "String Normal",
				data: "Hello, World!",
			},
			{
				name: "Bytes Normal",
				data: []byte("Hello, World!"),
			},
			{
				name: "String Chinese",
				data: "你好，世界！",
			},
			{
				name: "String With Key",
				data: "Secret Message",
				key:  "CustomKey",
			},

			// 空值测试
			{
				name: "String Empty",
				data: "",
			},
			{
				name: "Bytes Empty",
				data: []byte{},
			},
			{
				name: "Nil Bytes",
				data: []byte(nil),
			},

			// 特殊内容测试
			{
				name: "Special Characters",
				data: "!@#$%^&*()_+-=[]{}|;:,.<>?",
			},
			{
				name: "Long String",
				data: strings.Repeat("Long Message ", 100),
			},
			{
				name: "Binary Data",
				data: []byte{0x00, 0xFF, 0x7F, 0x80},
			},

			// 错误情况测试
			{
				name:     "Invalid Type",
				data:     123,
				wantFail: true,
			},
			{
				name: "Binary Key",
				data: "test",
				key:  string([]byte{0xFF, 0x00}), // 二进制密钥是合法的
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var encrypted string

				// 执行加密
				switch v := tt.data.(type) {
				case string:
					encrypted = Encrypt(v, tt.key)
					// 检查空输入
					if v == "" {
						assert.Empty(t, encrypted)
						return
					}
				case []byte:
					encrypted = Encrypt(v, tt.key)
					// 检查空输入
					if len(v) == 0 {
						assert.Empty(t, encrypted)
						return
					}
				default:
					if !tt.wantFail {
						t.Errorf("unexpected data type: %T", tt.data)
					}
					return
				}

				// 检查失败情况
				if tt.wantFail {
					assert.Empty(t, encrypted, "expected encryption to fail")
					return
				}

				// 验证加密结果
				assert.NotEmpty(t, encrypted, "encrypted string should not be empty for non-empty input")
				_, err := base64.StdEncoding.DecodeString(encrypted)
				assert.NoError(t, err, "encrypted string should be valid base64")

				// 执行解密
				switch expected := tt.data.(type) {
				case string:
					decrypted := Decrypt[string](encrypted, tt.key)
					assert.Equal(t, expected, decrypted)
				case []byte:
					decrypted := Decrypt[[]byte](encrypted, tt.key)
					if len(decrypted) == 0 && len(expected) == 0 {
						decrypted = []byte{} // 确保空字节数组的类型一致
					}
					assert.Equal(t, expected, decrypted)
				}
			})
		}
	})

	t.Run("Edge", func(t *testing.T) {
		tests := []struct {
			name string
			test func(t *testing.T)
		}{
			{
				name: "Wrong Key",
				test: func(t *testing.T) {
					data := "Secret Message"
					encrypted := Encrypt(data, "key1")
					decrypted := Decrypt[string](encrypted, "key2")
					assert.NotEqual(t, data, decrypted)
				},
			},
			{
				name: "Invalid Base64",
				test: func(t *testing.T) {
					result := Decrypt[string]("invalid base64")
					assert.Empty(t, result)
				},
			},
			{
				name: "Long Key",
				test: func(t *testing.T) {
					data := "Test Message"
					longKey := strings.Repeat("key", 10)
					encrypted := Encrypt(data, longKey)
					decrypted := Decrypt[string](encrypted, longKey)
					assert.Equal(t, data, decrypted)
				},
			},
			{
				name: "Short Key",
				test: func(t *testing.T) {
					data := "Test Message"
					shortKey := "k"
					encrypted := Encrypt(data, shortKey)
					decrypted := Decrypt[string](encrypted, shortKey)
					assert.Equal(t, data, decrypted)
				},
			},
			{
				name: "Empty Key",
				test: func(t *testing.T) {
					data := "Test Message"
					encrypted := Encrypt(data, "")
					decrypted := Decrypt[string](encrypted, "")
					assert.Equal(t, data, decrypted)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, tt.test)
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		const goroutines = 100
		done := make(chan bool, goroutines)
		var wg sync.WaitGroup
		wg.Add(goroutines)

		data := "Concurrent Test Message"
		key := "TestKey"

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				// 加密
				encrypted := Encrypt(data, key)
				assert.NotEmpty(t, encrypted)

				// 解密
				decrypted := Decrypt[string](encrypted, key)
				assert.Equal(t, data, decrypted)

				done <- true
			}()
		}

		go func() {
			wg.Wait()
			close(done)
		}()

		for range done {
			// 等待所有 goroutine 完成
		}
	})
}
