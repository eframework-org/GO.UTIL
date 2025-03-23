// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XPrefs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/illumitacit/gostd/quit"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	pf := New()

	t.Run("Set and Unset", func(t *testing.T) {
		// 测试设置和取消设置键值
		assert.False(t, pf.Has("nonexistent"))
		pf.Set("key", "value")
		assert.True(t, pf.Has("key"))

		pf.Unset("key")
		assert.False(t, pf.Has("key"))
	})

	t.Run("Set and Get", func(t *testing.T) {
		// 基本类型测试
		tests := []struct {
			name     string
			key      string
			value    any
			getFunc  func(string) any
			expected any
		}{
			{"String", "strKey", "value", func(k string) any { return pf.GetString(k) }, "value"},
			{"Int", "intKey", 42, func(k string) any { return pf.GetInt(k) }, 42},
			{"Bool", "boolKey", true, func(k string) any { return pf.GetBool(k) }, true},
			{"Float", "floatKey", 3.14, func(k string) any { return pf.GetFloat(k) }, float32(3.14)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pf.Set(tt.key, tt.value)
				assert.Equal(t, tt.expected, tt.getFunc(tt.key))
			})
		}

		// 切片类型测试
		sliceTests := []struct {
			name     string
			key      string
			value    any
			getFunc  func(string) any
			expected any
		}{
			{"String Slice", "strSlice", []string{"a", "b", "c"},
				func(k string) any { return pf.GetStrings(k) }, []string{"a", "b", "c"}},
			{"Int Slice", "intSlice", []int{1, 2, 3},
				func(k string) any { return pf.GetInts(k) }, []int{1, 2, 3}},
			{"Float Slice", "floatSlice", []float32{1.1, 2.2, 3.3},
				func(k string) any { return pf.GetFloats(k) }, []float32{1.1, 2.2, 3.3}},
			{"Bool Slice", "boolSlice", []bool{true, false, true},
				func(k string) any { return pf.GetBools(k) }, []bool{true, false, true}},
		}

		for _, tt := range sliceTests {
			t.Run(tt.name, func(t *testing.T) {
				pf.Set(tt.key, tt.value)
				assert.Equal(t, tt.expected, tt.getFunc(tt.key))
			})
		}
	})

	t.Run("Default Values", func(t *testing.T) {
		// 测试默认值
		assert.Equal(t, "default", pf.GetString("missing", "default"))
		assert.Equal(t, 100, pf.GetInt("missing", 100))
		assert.Equal(t, true, pf.GetBool("missing", true))
		assert.Equal(t, float32(1.23), pf.GetFloat("missing", 1.23))

		// 切片类型默认值
		assert.Equal(t, []string{"default"}, pf.GetStrings("missing", []string{"default"}))
		assert.Equal(t, []int{1, 2}, pf.GetInts("missing", []int{1, 2}))
		assert.Equal(t, []float32{1.1}, pf.GetFloats("missing", []float32{1.1}))
		assert.Equal(t, []bool{true}, pf.GetBools("missing", []bool{true}))
	})
}

func TestSources(t *testing.T) {
	// 初始化测试数据
	Asset().Set("intKey", 42)
	Asset().Set("intsKey", []any{1, 2, 3})
	Asset().Set("stringKey", "assetValue")
	Asset().Set("floatKey", 3.14)
	Asset().Set("boolKey", true)
	Asset().Set("stringsKey", []any{"a", "b", "c"})
	Asset().Set("floatsKey", []any{1.1, 2.2, 3.3})
	Asset().Set("boolsKey", []any{true, false, true})

	// 初始化本地测试数据
	Local().Set("localIntKey", 100)
	Local().Set("localIntsKey", []any{4, 5, 6})
	Local().Set("localStringKey", "localValue")
	Local().Set("overrideKey", "localOverride")

	t.Run("HasKey Tests", func(t *testing.T) {
		t.Run("Asset Only", func(t *testing.T) {
			// 测试仅在 Asset 中存在的键
			assert.True(t, HasKey("intKey"))
			assert.False(t, HasKey("nonexistentKey"))
		})

		t.Run("Multiple Sources", func(t *testing.T) {
			// 测试多个来源的键
			assert.True(t, HasKey("localIntKey", Local()))
			assert.True(t, HasKey("intKey", Local(), Asset()))
			assert.False(t, HasKey("nonexistentKey", Local(), Asset()))
		})
	})

	t.Run("GetInt Tests", func(t *testing.T) {
		t.Run("Asset Default", func(t *testing.T) {
			// 测试从 Asset 获取默认值
			assert.Equal(t, 42, GetInt("intKey", 0))
		})

		t.Run("Local Source", func(t *testing.T) {
			// 测试从 Local 获取值
			assert.Equal(t, 100, GetInt("localIntKey", 0, Local()))
		})

		t.Run("Default Value", func(t *testing.T) {
			// 测试获取不存在的键的默认值
			assert.Equal(t, 999, GetInt("nonexistentKey", 999, Local(), Asset()))
		})

		t.Run("Float to Int Conversion", func(t *testing.T) {
			// 测试浮点数转换为整数
			assert.Equal(t, 3, GetInt("floatKey", 0))
		})
	})

	t.Run("GetInts Tests", func(t *testing.T) {
		t.Run("Asset Default", func(t *testing.T) {
			// 测试从 Asset 获取整数切片的默认值
			expected := []int{1, 2, 3}
			assert.Equal(t, expected, GetInts("intsKey", nil))
		})

		t.Run("Local Source", func(t *testing.T) {
			// 测试从 Local 获取整数切片
			expected := []int{4, 5, 6}
			assert.Equal(t, expected, GetInts("localIntsKey", nil, Local()))
		})

		t.Run("Default Value", func(t *testing.T) {
			// 测试获取不存在的键的默认整数切片
			defaultVal := []int{7, 8, 9}
			assert.Equal(t, defaultVal, GetInts("nonexistentKey", defaultVal, Local(), Asset()))
		})
	})

	t.Run("Get Tests", func(t *testing.T) {
		t.Run("String Value", func(t *testing.T) {
			// 测试获取字符串值
			assert.Equal(t, "assetValue", Get("stringKey", ""))
		})

		t.Run("Float Value", func(t *testing.T) {
			// 测试获取浮点值
			assert.Equal(t, 3.14, Get("floatKey", 0.0))
		})

		t.Run("Bool Value", func(t *testing.T) {
			// 测试获取布尔值
			assert.Equal(t, true, Get("boolKey", false))
		})

		t.Run("Source Priority", func(t *testing.T) {
			// 测试源优先级
			assert.Equal(t, "localOverride", Get("overrideKey", "", Local(), Asset()))
		})
	})

	t.Run("Type Specific Tests", func(t *testing.T) {
		t.Run("GetString", func(t *testing.T) {
			// 测试获取字符串
			assert.Equal(t, "assetValue", GetString("stringKey", ""))
			assert.Equal(t, "default", GetString("nonexistentKey", "default"))
		})

		t.Run("GetStrings", func(t *testing.T) {
			// 测试获取字符串切片
			expected := []string{"a", "b", "c"}
			assert.Equal(t, expected, GetStrings("stringsKey", nil))
		})

		t.Run("GetFloat", func(t *testing.T) {
			// 测试获取浮点数
			assert.Equal(t, float32(3.14), GetFloat("floatKey", 0))
		})

		t.Run("GetFloats", func(t *testing.T) {
			// 测试获取浮点数切片
			expected := []float32{1.1, 2.2, 3.3}
			result := GetFloats("floatsKey", nil)
			for i, v := range expected {
				assert.InDelta(t, v, result[i], 0.001)
			}
		})

		t.Run("GetBool", func(t *testing.T) {
			// 测试获取布尔值
			assert.True(t, GetBool("boolKey", false))
		})

		t.Run("GetBools", func(t *testing.T) {
			// 测试获取布尔值切片
			expected := []bool{true, false, true}
			assert.Equal(t, expected, GetBools("boolsKey", nil))
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("Nil Sources", func(t *testing.T) {
			// 测试空源
			assert.Equal(t, 42, GetInt("intKey", 0, nil))
		})

		t.Run("Empty Sources", func(t *testing.T) {
			// 测试空源
			assert.Equal(t, 42, GetInt("intKey", 0))
		})

		t.Run("Type Mismatches", func(t *testing.T) {
			// 测试类型不匹配
			Asset().Set("mismatchKey", "not an int")
			assert.Equal(t, 0, GetInt("mismatchKey"))
		})
	})
}

func TestAutoSave(t *testing.T) {
	tmpDir := t.TempDir()
	localDir := filepath.Join(tmpDir, "Local")
	err := os.MkdirAll(localDir, 0755)
	assert.NoError(t, err)

	localFile := filepath.Join(localDir, "Preferences.json")
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	reset()
	os.Args = []string{"test", "--Prefs@Local=" + localFile}

	Local().Set("test_key", "test_value")

	quit.GetQuitChannel() <- struct{}{}
	quit.GetWaiter().Wait()

	assert.FileExists(t, localFile)
	Local().read(localFile)
	assert.Equal(t, "test_value", Local().GetString("test_key"))
}

func TestRead(t *testing.T) {
	pf := &prefsLocal{}

	// 创建临时目录
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_prefs.json")

	// 准备测试数据
	data := `{
		"stringKey": "stringValue",
		"intKey": 123,
		"boolKey": true,
		"intSliceKey": [1, 2, 3],
		"floatSliceKey": [1.1, 2.2, 3.3],
		"stringSliceKey": ["a", "b", "c"],
		"boolSliceKey": [true, false, true]
	}`

	// 写入测试文件
	err := os.MkdirAll(filepath.Dir(testFile), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(testFile, []byte(data), 0644)
	assert.NoError(t, err)

	// 测试读取偏好设置
	t.Run("Read Preferences", func(t *testing.T) {
		assert.True(t, pf.read(testFile), "Should read file successfully")

		// 验证各种类型的数据
		assert.Equal(t, "stringValue", pf.GetString("stringKey"))
		assert.Equal(t, 123, pf.GetInt("intKey"))
		assert.True(t, pf.GetBool("boolKey"))
		assert.Equal(t, []int{1, 2, 3}, pf.GetInts("intSliceKey"))
		assert.Equal(t, []float32{1.1, 2.2, 3.3}, pf.GetFloats("floatSliceKey"))
		assert.Equal(t, []string{"a", "b", "c"}, pf.GetStrings("stringSliceKey"))
		assert.Equal(t, []bool{true, false, true}, pf.GetBools("boolSliceKey"))
	})

	// 测试保存和重新读取
	t.Run("Save and Read Again", func(t *testing.T) {
		pf.file = testFile
		assert.True(t, pf.Save(), "Should save file successfully")

		// 清空当前数据
		pf = &prefsLocal{}

		// 重新读取并验证
		assert.True(t, pf.read(testFile), "Should read file again successfully")
		assert.Equal(t, "stringValue", pf.GetString("stringKey"))
		assert.Equal(t, 123, pf.GetInt("intKey"))
	})

	// 测试读取不存在的文件
	t.Run("Read Non-existent File", func(t *testing.T) {
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.json")
		assert.False(t, pf.read(nonExistentFile), "Should fail reading non-existent file")
	})

	// 测试读取无效的 JSON
	t.Run("Read Invalid JSON", func(t *testing.T) {
		invalidFile := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(invalidFile, []byte("invalid json"), 0644)
		assert.NoError(t, err)
		assert.False(t, pf.read(invalidFile), "Should fail reading invalid JSON")
	})

	t.Run("Complex JSON", func(t *testing.T) {
		complexData := `{
			"nullValue": null,
			"emptyObject": {},
			"emptyArray": [],
			"nestedObject": {
				"key": "value"
			},
			"mixedArray": [1, "two", true, null]
		}`

		complexFile := filepath.Join(tmpDir, "complex.json")
		err := os.WriteFile(complexFile, []byte(complexData), 0644)
		assert.NoError(t, err)

		pf := &prefsLocal{}
		assert.True(t, pf.read(complexFile))

		assert.Nil(t, pf.Get("nullValue"))
		assert.NotNil(t, pf.Get("emptyObject"))
		assert.NotNil(t, pf.Get("emptyArray"))
		assert.NotNil(t, pf.Get("nestedObject"))
		assert.NotNil(t, pf.Get("mixedArray"))
	})

	t.Run("Large File", func(t *testing.T) {
		// 生成大文件
		largeData := make(map[string]string)
		for i := 0; i < 1000; i++ {
			largeData[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
		}

		jsonData, err := json.Marshal(largeData)
		assert.NoError(t, err)

		largeFile := filepath.Join(tmpDir, "large.json")
		err = os.WriteFile(largeFile, jsonData, 0644)
		assert.NoError(t, err)

		pf := &prefsLocal{}
		assert.True(t, pf.read(largeFile))
		assert.Equal(t, "value42", pf.GetString("key42"))
	})
}

func TestEval(t *testing.T) {
	t.Run("Basic Replacement", func(t *testing.T) {
		pf := New()
		pf.Set("name", "John")
		pf.Set("greeting", "Hello ${Prefs.name}")

		result := pf.Eval("${Prefs.greeting}")
		assert.Equal(t, "Hello John", result)
	})

	t.Run("Missing Variable", func(t *testing.T) {
		pf := New()
		result := pf.Eval("${Prefs.missing}")
		assert.Equal(t, "${Prefs.missing}(Unknown)", result)
	})

	t.Run("Recursive Variables", func(t *testing.T) {
		pf := New()
		pf.Set("recursive1", "${Prefs.recursive2}")
		pf.Set("recursive2", "${Prefs.recursive1}")

		result := pf.Eval("${Prefs.recursive1}")
		assert.Equal(t, "${Prefs.recursive1}(Recursive)", result)
	})

	t.Run("Nested Variables", func(t *testing.T) {
		pf := New()
		pf.Set("outer", "value")

		result := pf.Eval("${Prefs.outer${Prefs.inner}}")
		assert.Equal(t, "${Prefs.outer${Prefs.inner}(Nested)}", result)
	})

	t.Run("Multiple Replacements", func(t *testing.T) {
		pf := New()
		pf.Set("first", "John")
		pf.Set("last", "Doe")
		pf.Set("missing", "")

		result := pf.Eval("${Prefs.first} ${Prefs.last} ${Prefs.missing} ${Prefs.unknown}")
		assert.Equal(t, "John Doe ${Prefs.missing}(Unknown) ${Prefs.unknown}(Unknown)", result)
	})

	t.Run("Empty Value", func(t *testing.T) {
		pf := New()
		pf.Set("empty", "")

		result := pf.Eval("test${Prefs.empty}end")
		assert.Equal(t, "test${Prefs.empty}(Unknown)end", result)
	})

	t.Run("Nested Fields", func(t *testing.T) {
		pf := New()
		pf.Set("first", "John")
		pf.Set("last", "Doe")

		// 创建并设置子对象
		child := New()
		child.Set("name", "Mike")
		pf.Set("child", child)

		result := pf.Eval("${Prefs.first} and ${Prefs.last} has a child named ${Prefs.child.name} age ${Prefs.child.age}")
		assert.Equal(t, "John and Doe has a child named Mike age ${Prefs.child.age}(Unknown)", result)
	})

	t.Run("Deep Nested Fields", func(t *testing.T) {
		pf := New()

		// 创建多层嵌套结构
		child := New()
		child.Set("name", "Mike")

		info := New()
		info.Set("age", "10")
		child.Set("info", info)

		pf.Set("child", child)

		// 测试多层嵌套访问
		result := pf.Eval("Child ${Prefs.child.name} is ${Prefs.child.info.age} years old")
		assert.Equal(t, "Child Mike is 10 years old", result)

		// 测试无效的多层路径
		result = pf.Eval("${Prefs.child.info.unknown}")
		assert.Equal(t, "${Prefs.child.info.unknown}(Unknown)", result)

		// 测试部分无效的路径
		result = pf.Eval("${Prefs.child.unknown.field}")
		assert.Equal(t, "${Prefs.child.unknown.field}(Unknown)", result)
	})
}

func TestOverrides(t *testing.T) {
	// 保存原始参数并在测试后恢复
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 创建临时文件
	tmpDir := t.TempDir()
	assetFile := filepath.Join(tmpDir, "Assets/Preferences.json")
	localFile := filepath.Join(tmpDir, "Local/Preferences.json")

	// 写入测试数据
	assetData := `{
		"name": "DefaultName",
		"age": 25,
		"setting": "default"
	}`
	localData := `{
		"name": "LocalName",
		"location": "LocalCity"
	}`

	err := os.MkdirAll(filepath.Dir(assetFile), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(assetFile, []byte(assetData), 0644)
	assert.NoError(t, err)

	err = os.MkdirAll(filepath.Dir(localFile), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(localFile, []byte(localData), 0644)
	assert.NoError(t, err)

	t.Run("Override Config Files", func(t *testing.T) {
		reset()
		// 设置命令行参数来覆盖配置文件路径
		os.Args = []string{"test",
			"--Prefs@Asset=" + assetFile,
			"--Prefs@Local=" + localFile,
		}
		// 验证文件被正确加载
		assert.Equal(t, "DefaultName", Asset().GetString("name"))
		assert.Equal(t, "LocalName", Local().GetString("name"))
	})

	t.Run("Override Config Values", func(t *testing.T) {
		reset()
		// 设置命令行参数来覆盖具体配置项
		os.Args = []string{"test",
			"--Prefs@Asset=" + assetFile,
			"--Prefs@Local=" + localFile,
			"--Prefs@Asset.name=OverriddenName",
			"--Prefs@Local.location=OverriddenCity",
			"--Prefs.Log.Std.Config.Level=Debug",
			"--Prefs@Asset.UI.Window.Style.Theme=Dark",
			"--Prefs@Local.Network.Server.Config.Port=8080",
		}
		// 验证配置项被正确覆盖
		assert.Equal(t, "OverriddenName", Asset().GetString("name"))
		assert.Equal(t, "OverriddenCity", Local().GetString("location"))
		// 验证未覆盖的配置项保持原值
		assert.Equal(t, float32(25), Asset().GetFloat("age"))
		assert.Equal(t, "Debug", Asset().Get("Log").(IBase).Get("Std").(IBase).Get("Config").(IBase).Get("Level"))
		assert.Equal(t, "Dark", Asset().Get("UI").(IBase).Get("Window").(IBase).Get("Style").(IBase).Get("Theme"))
		assert.Equal(t, 8080, Local().Get("Network").(IBase).Get("Server").(IBase).Get("Config").(IBase).GetInt("Port"))
	})

	t.Run("Invalid Config Files", func(t *testing.T) {
		reset()
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.json")
		os.Args = []string{"test",
			"--Prefs@Asset=" + nonExistentFile,
		}
		// 验证处理不存在的文件
		assert.False(t, Asset().Has("name"))
	})

	t.Run("Mixed Overrides", func(t *testing.T) {
		reset()
		os.Args = []string{"test",
			"--Prefs@Asset=" + assetFile,
			"--Prefs@Local=" + localFile,
			"--Prefs@Asset.setting=overridden",
			"--Prefs@Local.newKey=newValue",
		}
		// 验证混合覆盖场景
		assert.Equal(t, "DefaultName", Asset().GetString("name"))
		assert.Equal(t, "overridden", Asset().GetString("setting"))
		assert.Equal(t, "newValue", Local().GetString("newKey"))
	})

	t.Run("Empty Values", func(t *testing.T) {
		reset()
		os.Args = []string{"test",
			"--Prefs@Asset=" + assetFile,
			"--Prefs@Asset.emptyKey=",
		}
		assert.Equal(t, "", Asset().GetString("emptyKey"))
	})

	t.Run("Invalid Key Format", func(t *testing.T) {
		reset()
		os.Args = []string{"test",
			"--Prefs@Asset=" + assetFile,
			"--Prefs@Invalid.key=value", // 错误的前缀
			"--Prefs@Asset.validKey=value",
		}
		assert.Equal(t, "value", Asset().GetString("validKey"))
		assert.False(t, Asset().Has("key"))
	})

	t.Run("Duplicate Keys", func(t *testing.T) {
		reset()
		os.Args = []string{"test",
			"--Prefs@Asset=" + assetFile,
			"--Prefs@Asset.key=first",
			"--Prefs@Asset.key=second", // 相同的键
		}
		assert.Equal(t, "second", Asset().GetString("key"))
	})
}
