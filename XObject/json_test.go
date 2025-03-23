// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XObject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用的结构体
type testStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Valid bool   `json:"valid"`
}

// TestToJson 测试对象到 JSON 的转换。
func TestToJson(t *testing.T) {
	t.Run("Basic Types", func(t *testing.T) {
		// 测试基本类型
		tests := []struct {
			name     string
			input    any
			expected string
		}{
			{"string", "hello", "hello"},
			{"int", 42, "42"},
			{"bool", true, "true"},
			{"float", 3.14, "3.14"},
			{"null", nil, "null"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result, err := ToJson(test.input)
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("Struct", func(t *testing.T) {
		obj := testStruct{
			Name:  "test",
			Age:   25,
			Valid: true,
		}
		result, err := ToJson(obj)
		assert.NoError(t, err)
		assert.Equal(t, `{"name":"test","age":25,"valid":true}`, result)
	})

	t.Run("Pretty Print", func(t *testing.T) {
		obj := testStruct{
			Name:  "test",
			Age:   25,
			Valid: true,
		}
		result, err := ToJson(obj, true)
		assert.NoError(t, err)
		expected := `{
    "name": "test",
    "age": 25,
    "valid": true
}`
		assert.Equal(t, expected, result)
	})

	t.Run("Invalid Types", func(t *testing.T) {
		// 测试无法序列化的类型
		ch := make(chan int)
		_, err := ToJson(ch)
		assert.Error(t, err)
	})
}

// TestFromJson 测试 JSON 到对象的转换。
func TestFromJson(t *testing.T) {
	t.Run("Basic Types", func(t *testing.T) {
		// 字符串到基本类型
		var strVal string
		assert.NoError(t, FromJson(`"hello"`, &strVal))
		assert.Equal(t, "hello", strVal)

		var intVal int
		assert.NoError(t, FromJson("42", &intVal))
		assert.Equal(t, 42, intVal)

		var boolVal bool
		assert.NoError(t, FromJson("true", &boolVal))
		assert.True(t, boolVal)

		var floatVal float64
		assert.NoError(t, FromJson("3.14", &floatVal))
		assert.Equal(t, 3.14, floatVal)
	})

	t.Run("Struct", func(t *testing.T) {
		jsonStr := `{"name":"test","age":25,"valid":true}`
		var obj testStruct
		err := FromJson(jsonStr, &obj)
		assert.NoError(t, err)
		assert.Equal(t, "test", obj.Name)
		assert.Equal(t, 25, obj.Age)
		assert.True(t, obj.Valid)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		var obj testStruct
		err := FromJson(`{"invalid json`, &obj)
		assert.Error(t, err)
	})

	t.Run("Type Mismatch", func(t *testing.T) {
		var obj testStruct
		err := FromJson(`{"name":42}`, &obj)
		assert.Error(t, err)
	})
}

// TestToByte 测试对象到字节数组的转换。
func TestToByte(t *testing.T) {
	t.Run("Basic Types", func(t *testing.T) {
		result, err := ToByte("hello")
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello"), result)

		result, err = ToByte(42)
		assert.NoError(t, err)
		assert.Equal(t, []byte("42"), result)
	})

	t.Run("Struct", func(t *testing.T) {
		obj := testStruct{
			Name:  "test",
			Age:   25,
			Valid: true,
		}
		result, err := ToByte(obj)
		assert.NoError(t, err)
		assert.Equal(t, []byte(`{"name":"test","age":25,"valid":true}`), result)
	})

	t.Run("Invalid Types", func(t *testing.T) {
		ch := make(chan int)
		_, err := ToByte(ch)
		assert.Error(t, err)
	})
}

// TestFromByte 测试字节数组到对象的转换。
func TestFromByte(t *testing.T) {
	t.Run("Basic Types", func(t *testing.T) {
		var strVal string
		assert.NoError(t, FromByte([]byte(`"hello"`), &strVal))
		assert.Equal(t, "hello", strVal)

		var intVal int
		assert.NoError(t, FromByte([]byte("42"), &intVal))
		assert.Equal(t, 42, intVal)
	})

	t.Run("Struct", func(t *testing.T) {
		jsonBytes := []byte(`{"name":"test","age":25,"valid":true}`)
		var obj testStruct
		err := FromByte(jsonBytes, &obj)
		assert.NoError(t, err)
		assert.Equal(t, "test", obj.Name)
		assert.Equal(t, 25, obj.Age)
		assert.True(t, obj.Valid)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		var obj testStruct
		err := FromByte([]byte(`{"invalid json`), &obj)
		assert.Error(t, err)
	})

	t.Run("Empty Input", func(t *testing.T) {
		var obj testStruct
		err := FromByte(nil, &obj)
		assert.Error(t, err)
	})
}
