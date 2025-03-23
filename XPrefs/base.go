// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XPrefs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/eframework-org/GO.UTIL/XObject"
)

// IBase 定义了配置管理的基础接口。
// 提供了配置项的增删改查、类型转换和格式化输出等功能。
type IBase interface {
	// Keys 返回所有配置项的键名列表。
	Keys() []string

	// Has 检查配置项是否存在。
	// 输入键名作为查询条件，如果配置项存在则返回 true，否则返回 false。
	Has(key string) bool

	// Set 设置配置项的值。
	// 输入键名和任意类型的值，设置完成后返回接口实例本身，支持链式调用。
	Set(key string, value any) IBase

	// Unset 删除指定的配置项。
	// 输入要删除的键名，删除完成后返回接口实例本身，支持链式调用。
	Unset(key string) IBase

	// Get 获取配置项的值。
	// 输入键名和可选的默认值，返回配置项的值，如果配置项不存在且未提供默认值则返回 nil。
	Get(key string, defval ...any) any

	// Gets 获取配置项的值数组。
	// 输入键名和可选的默认值数组，返回配置项的值数组，如果配置项不存在且未提供默认值则返回 nil。
	Gets(key string, defval ...[]any) []any

	// GetInt 获取配置项的整数值。
	// 输入键名和可选的默认值，返回配置项的整数值，如果配置项不存在或类型转换失败且未提供默认值则返回 0。
	GetInt(key string, defval ...int) int

	// GetInts 获取配置项的整数数组。
	// 输入键名和可选的默认值数组，返回配置项的整数数组，如果配置项不存在或类型转换失败且未提供默认值则返回 nil。
	GetInts(key string, defval ...[]int) []int

	// GetFloat 获取配置项的浮点数值。
	// 输入键名和可选的默认值，返回配置项的浮点数值，如果配置项不存在或类型转换失败且未提供默认值则返回 0。
	GetFloat(key string, defval ...float32) float32

	// GetFloats 获取配置项的浮点数数组。
	// 输入键名和可选的默认值数组，返回配置项的浮点数数组，如果配置项不存在或类型转换失败且未提供默认值则返回 nil。
	GetFloats(key string, defval ...[]float32) []float32

	// GetBool 获取配置项的布尔值。
	// 输入键名和可选的默认值，返回配置项的布尔值，如果配置项不存在或类型转换失败且未提供默认值则返回 false。
	GetBool(key string, defval ...bool) bool

	// GetBools 获取配置项的布尔值数组。
	// 输入键名和可选的默认值数组，返回配置项的布尔值数组，如果配置项不存在或类型转换失败且未提供默认值则返回 nil。
	GetBools(key string, defval ...[]bool) []bool

	// GetString 获取配置项的字符串值。
	// 输入键名和可选的默认值，返回配置项的字符串值，如果配置项不存在或类型转换失败且未提供默认值则返回空字符串。
	GetString(key string, defval ...string) string

	// GetStrings 获取配置项的字符串数组。
	// 输入键名和可选的默认值数组，返回配置项的字符串数组，如果配置项不存在或类型转换失败且未提供默认值则返回 nil。
	GetStrings(key string, defval ...[]string) []string

	// Json 将配置内容转换为 JSON 字符串。
	// 输入可选的格式化参数，如果为 true 则返回格式化后的 JSON 字符串，默认返回压缩的 JSON 字符串。
	Json(pretty ...bool) string

	// Eval 计算包含配置项引用的字符串表达式。
	// 输入包含配置项引用的字符串表达式，返回计算后的结果字符串。
	Eval(input string) string
}

// prefsBase 实现了 IBase 接口，提供配置管理的基础功能。
// 支持并发安全的配置项读写、类型转换、多级配置和变量引用等功能。
type prefsBase struct {
	sync.RWMutex                // 读写锁，用于保证并发操作的安全性
	pairs        map[string]any // 存储配置项的键值对映射
	npairs       map[string]any // 存储多级配置的缓存，避免重复解析
	keys         []string       // 存储所有配置项的键名，用于快速遍历
}

// Keys 返回所有配置项的键名列表。
// 首次调用时会构建键名缓存，后续调用直接返回缓存的副本。
// 返回的切片是内部键名列表的一个副本，可以安全地修改。
func (pb *prefsBase) Keys() []string {
	pb.RLock()
	if pb.keys != nil {
		keys := make([]string, len(pb.keys))
		copy(keys, pb.keys)
		pb.RUnlock()
		return keys
	}
	pb.RUnlock()

	pb.Lock()
	defer pb.Unlock()

	pb.keys = make([]string, 0, len(pb.pairs))
	for key := range pb.pairs {
		pb.keys = append(pb.keys, key)
	}
	keys := make([]string, len(pb.keys))
	copy(keys, pb.keys)
	return keys
}

// Has 检查配置项是否存在。
// 输入键名作为查询条件，在配置项映射中查找对应的键。
// 返回 true 表示配置项存在，false 表示不存在。
// 此方法是并发安全的，使用读锁保护。
func (pb *prefsBase) Has(key string) bool {
	pb.RLock()
	defer pb.RUnlock()

	_, exists := pb.pairs[key]
	return exists
}

// Set 设置配置项的值。
// 输入键名和任意类型的值，将其存储在配置项映射中。
// 如果键不存在，会将其添加到键名列表中。
// 返回接口实例本身，支持链式调用。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) Set(key string, value any) IBase {
	pb.Lock()
	defer pb.Unlock()

	if pb.pairs == nil {
		pb.pairs = make(map[string]any, 8)
		pb.keys = make([]string, 0, 8)
	}

	if _, exists := pb.pairs[key]; !exists {
		pb.keys = append(pb.keys, key)
	}
	pb.pairs[key] = value
	return pb
}

// Unset 删除指定的配置项。
// 输入要删除的键名，会同时从配置项映射和键名列表中移除。
// 如果键存在于多级配置缓存中，也会一并删除。
// 返回接口实例本身，支持链式调用。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) Unset(key string) IBase {
	pb.Lock()
	defer pb.Unlock()

	if pb.npairs != nil {
		if _, exists := pb.npairs[key]; exists {
			delete(pb.npairs, key)
		}
	}

	if _, exists := pb.pairs[key]; exists {
		delete(pb.pairs, key)
		for i, k := range pb.keys {
			if k == key {
				pb.keys = append(pb.keys[:i], pb.keys[i+1:]...)
				break
			}
		}
	}
	return pb
}

// Get 获取配置项的值。
// 输入键名和可选的默认值。首先检查多级配置缓存，然后查找配置项映射。
// 如果值是映射类型，会创建一个新的配置实例并缓存。
// 如果键不存在且提供了默认值，则返回默认值。
// 此方法是并发安全的，使用读锁保护。
func (pb *prefsBase) Get(key string, defval ...any) any {
	pb.RLock()
	defer pb.RUnlock()

	if pb.npairs != nil {
		if val, exists := pb.npairs[key]; exists {
			return val
		}
	}

	if val, exists := pb.pairs[key]; exists {
		switch val.(type) {
		case map[string]any:
			nv := &prefsBase{pairs: val.(map[string]any)}
			if pb.npairs == nil {
				pb.npairs = make(map[string]any)
			}
			pb.npairs[key] = nv
			return nv
		default:
			return val
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return nil
}

// Gets 获取配置项的值数组。
// 输入键名和可选的默认值数组。首先检查多级配置缓存，然后查找配置项映射。
// 支持处理映射数组，会将每个映射转换为配置实例并缓存。
// 如果键不存在且提供了默认值数组，则返回默认值数组。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) Gets(key string, defval ...[]any) []any {
	pb.Lock()
	defer pb.Unlock()

	if pb.npairs != nil {
		if val, exists := pb.npairs[key]; exists {
			switch v := val.(type) {
			case []any:
				return v
			}
		}
	}

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case []map[string]any:
			var nv []any
			for _, vv := range v {
				nv = append(nv, &prefsBase{pairs: vv})
			}
			if pb.npairs == nil {
				pb.npairs = make(map[string]any)
			}
			pb.npairs[key] = nv
			return nv
		case []any:
			return v
		}
	}

	if len(defval) > 0 {
		return defval[0]
	}

	return nil
}

// GetInt 获取配置项的整数值。
// 输入键名和可选的默认值。支持多种数值类型的自动转换，包括：
// - 各种整数类型（int8 到 int64，uint8 到 uint64）
// - 浮点数类型（float32，float64）
// - 字符串类型（会尝试解析为整数）
// 如果转换失败且提供了默认值，则返回默认值；否则返回 0。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetInt(key string, defval ...int) int {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case int:
			return v
		case int8:
			return int(v)
		case int16:
			return int(v)
		case int32:
			return int(v)
		case int64:
			return int(v)
		case uint:
			return int(v)
		case uint8:
			return int(v)
		case uint16:
			return int(v)
		case uint32:
			return int(v)
		case uint64:
			return int(v)
		case float32:
			return int(v)
		case float64:
			return int(v)
		case string:
			if intVal, err := strconv.Atoi(v); err == nil {
				return intVal
			}
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return 0
}

// GetInts 获取配置项的整数数组。
// 输入键名和可选的默认值数组。支持处理原生整数数组和任意类型数组。
// 对于任意类型数组，会尝试将每个元素转换为整数，支持与 GetInt 相同的类型转换。
// 如果键不存在或转换失败且提供了默认值数组，则返回默认值数组。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetInts(key string, defval ...[]int) []int {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case []int:
			return v
		case []any:
			intSlice := make([]int, len(v))
			for i, vv := range v {
				switch ve := vv.(type) {
				case int:
					intSlice[i] = ve
				case int8:
					intSlice[i] = int(ve)
				case int16:
					intSlice[i] = int(ve)
				case int32:
					intSlice[i] = int(ve)
				case int64:
					intSlice[i] = int(ve)
				case uint:
					intSlice[i] = int(ve)
				case uint8:
					intSlice[i] = int(ve)
				case uint16:
					intSlice[i] = int(ve)
				case uint32:
					intSlice[i] = int(ve)
				case uint64:
					intSlice[i] = int(ve)
				case float32:
					intSlice[i] = int(ve)
				case float64:
					intSlice[i] = int(ve)
				case string:
					if iv, err := strconv.Atoi(ve); err == nil {
						intSlice[i] = iv
					}
				}
			}
			return intSlice
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return nil
}

// GetFloat 获取配置项的浮点数值。
// 输入键名和可选的默认值。支持多种数值类型的自动转换，包括：
// - 各种整数类型（int8 到 int64，uint8 到 uint64）
// - 浮点数类型（float32，float64）
// - 字符串类型（会尝试解析为浮点数）
// 如果转换失败且提供了默认值，则返回默认值；否则返回 0。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetFloat(key string, defval ...float32) float32 {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case int:
			return float32(v)
		case int8:
			return float32(v)
		case int16:
			return float32(v)
		case int32:
			return float32(v)
		case int64:
			return float32(v)
		case uint:
			return float32(v)
		case uint8:
			return float32(v)
		case uint16:
			return float32(v)
		case uint32:
			return float32(v)
		case uint64:
			return float32(v)
		case float32:
			return float32(v)
		case float64:
			return float32(v)
		case string:
			if nval, err := strconv.ParseFloat(v, 32); err == nil {
				return float32(nval)
			}
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return 0
}

// GetFloats 获取配置项的浮点数数组。
// 输入键名和可选的默认值数组。支持处理原生浮点数数组和任意类型数组。
// 对于任意类型数组，会尝试将每个元素转换为浮点数。
// 如果键不存在或转换失败且提供了默认值数组，则返回默认值数组。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetFloats(key string, defval ...[]float32) []float32 {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case []float32:
			return v
		case []any:
			floatSlice := make([]float32, len(v))
			for i, item := range v {
				if floatVal, ok := item.(float64); ok { // JSON numbers are float64
					floatSlice[i] = float32(floatVal)
				}
			}
			return floatSlice
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return nil
}

// GetBool 获取配置项的布尔值。
// 输入键名和可选的默认值。仅支持原生布尔类型的值。
// 如果值不是布尔类型且提供了默认值，则返回默认值；否则返回 false。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetBool(key string, defval ...bool) bool {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return false
}

// GetBools 获取配置项的布尔值数组。
// 输入键名和可选的默认值数组。支持处理原生布尔值数组和任意类型数组。
// 对于任意类型数组，会尝试将每个元素转换为布尔值。
// 如果键不存在或转换失败且提供了默认值数组，则返回默认值数组。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetBools(key string, defval ...[]bool) []bool {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case []bool:
			return v
		case []any:
			boolSlice := make([]bool, len(v))
			for i, item := range v {
				if boolVal, ok := item.(bool); ok {
					boolSlice[i] = boolVal
				}
			}
			return boolSlice
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return nil
}

// GetString 获取配置项的字符串值。
// 输入键名和可选的默认值。仅支持原生字符串类型的值。
// 如果值不是字符串类型且提供了默认值，则返回默认值；否则返回空字符串。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetString(key string, defval ...string) string {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return ""
}

// GetStrings 获取配置项的字符串数组。
// 输入键名和可选的默认值数组。支持处理原生字符串数组和任意类型数组。
// 对于任意类型数组，会尝试将每个元素转换为字符串。
// 如果键不存在或转换失败且提供了默认值数组，则返回默认值数组。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) GetStrings(key string, defval ...[]string) []string {
	pb.Lock()
	defer pb.Unlock()

	if val, exists := pb.pairs[key]; exists {
		switch v := val.(type) {
		case []string:
			return v
		case []any:
			strSlice := make([]string, len(v))
			for i, item := range v {
				if strVal, ok := item.(string); ok {
					strSlice[i] = strVal
				}
			}
			return strSlice
		}
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return nil
}

// Json 将配置内容转换为 JSON 字符串。
// 输入可选的格式化参数，如果为 true 则返回格式化后的 JSON 字符串。
// 如果未提供格式化参数或为 false，则返回压缩的 JSON 字符串。
// 此方法是并发安全的，使用写锁保护。
func (pb *prefsBase) Json(pretty ...bool) string {
	pb.Lock()
	defer pb.Unlock()

	ret, _ := XObject.ToJson(pb.pairs, pretty...)
	return ret
}

// Eval 计算包含配置项引用的字符串表达式。
// 输入包含配置项引用的字符串表达式，格式为 ${Prefs.key}。
// 支持以下特性：
// - 多级配置引用，如 ${Prefs.UI.Theme}
// - 循环引用检测，避免死循环
// - 嵌套变量引用检测
// - 未定义变量和空值处理
// 返回计算后的结果字符串，对于特殊情况会添加相应的标记。
func (pb *prefsBase) Eval(input string) string {
	pattern := regexp.MustCompile(`\$\{Prefs\.([^}]+?)\}`)
	visited := make(map[string]bool)

	var replaceFunc func(string) string
	replaceFunc = func(match string) string {
		path := pattern.FindStringSubmatch(match)[1]

		// 1. 检查嵌套变量
		if strings.Contains(path, "${") {
			return fmt.Sprintf("%v(Nested)", match)
		}

		// 2. 检查循环引用
		if visited[path] {
			return fmt.Sprintf("${Prefs.%v}(Recursive)", path)
		}
		visited[path] = true
		defer delete(visited, path)

		// 3. 获取变量值（支持多级路径）
		var value string
		if strings.Contains(path, ".") {
			parts := strings.Split(path, ".")
			current := (any)(pb).(IBase)
			for i := 0; i < len(parts)-1; i++ {
				if !current.Has(parts[i]) {
					return fmt.Sprintf("${Prefs.%v}(Unknown)", path)
				}
				next := current.Get(parts[i])
				if next == nil {
					return fmt.Sprintf("${Prefs.%v}(Unknown)", path)
				}
				if base, ok := next.(IBase); ok {
					current = base
				} else {
					return fmt.Sprintf("${Prefs.%v}(Unknown)", path)
				}
			}
			value = current.GetString(parts[len(parts)-1])
		} else {
			if !pb.Has(path) {
				return fmt.Sprintf("${Prefs.%v}(Unknown)", path)
			}
			value = pb.GetString(path)
		}

		// 4. 检查空值
		if value == "" {
			return fmt.Sprintf("${Prefs.%v}(Unknown)", path)
		}

		// 5. 递归处理值中的变量
		return pattern.ReplaceAllStringFunc(value, replaceFunc)
	}

	return pattern.ReplaceAllStringFunc(input, replaceFunc)
}

// parse 解析配置数据并应用命令行参数覆盖。
// 输入字节数组形式的配置数据，将其解析为配置项映射。
// 解析完成后会检查命令行参数中以 "Prefs." 开头的配置项，并用其值覆盖相应的配置。
// 支持多级配置的命令行参数覆盖。
// 返回 true 表示解析成功，false 表示解析失败。
func (pb *prefsBase) parse(data []byte) bool {
	defer func() {
		args := parseArgs()
		for k, v := range args {
			if strings.HasPrefix(k, "Prefs.") {
				key := strings.TrimPrefix(k, "Prefs.")
				if strings.Contains(key, ".") {
					parts := strings.Split(key, ".")
					current := (any)(pb).(IBase)
					for i := 0; i < len(parts)-1; i++ {
						part := parts[i]
						if !current.Has(part) {
							current.Set(part, New())
						}
						current = current.Get(part).(IBase)
					}
					current.Set(parts[len(parts)-1], v)
				} else {
					pb.Set(key, v)
				}
				fmt.Printf("XPrefs.Base.Parse: override %s = %s\n", key, v)
			}
		}
	}()

	if data == nil || len(data) == 0 {
		fmt.Printf("XPrefs.Base.Parse: nil data\n")
		return false
	}
	var kvs map[string]any
	err := json.Unmarshal(data, &kvs)
	if err != nil {
		fmt.Printf("XPrefs.Base.Parse: unmarshal error: %v\n", err)
		return false
	}
	for k, v := range kvs {
		pb.Set(k, v)
	}
	return true
}
