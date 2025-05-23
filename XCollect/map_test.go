// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XCollect

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

// 测试 Map 模块
func TestMap(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		m := NewMap()

		// 测试加载不存在的键
		value, _ := m.Load("不存在的键")
		if value != nil {
			t.Fatalf("Load 不存在的键应返回 nil，实际返回 %v", value)
		}

		// 测试加载存在的键
		m.Store("测试键", "测试值")
		value, _ = m.Load("测试键")
		if value != "测试值" {
			t.Fatalf("Load 存在的键应返回对应值，期望 %v，实际返回 %v", "测试值", value)
		}

		// 测试大量数据
		for i := range 1000 {
			key := fmt.Sprintf("key-%d", i)
			val := fmt.Sprintf("value-%d", i)
			m.Store(key, val)
		}

		for i := range 1000 {
			key := fmt.Sprintf("key-%d", i)
			expectedVal := fmt.Sprintf("value-%d", i)
			val, _ := m.Load(key)
			if val != expectedVal {
				t.Fatalf("Load 大量数据测试失败，键 %s 期望值 %s，实际值 %v", key, expectedVal, val)
			}
		}
	})

	t.Run("Store", func(t *testing.T) {
		m := NewMap()

		// 测试存储新键值对
		m.Store("key1", "value1")
		value, _ := m.Load("key1")
		if value != "value1" {
			t.Fatalf("Store 后 Load 应返回存储的值，期望 %v，实际返回 %v", "value1", value)
		}

		// 测试更新现有键值对
		m.Store("key1", "value1-updated")
		value, _ = m.Load("key1")
		if value != "value1-updated" {
			t.Fatalf("更新后 Load 应返回更新的值，期望 %v，实际返回 %v", "value1-updated", value)
		}

		// 测试存储不同类型的键
		m.Store(123, "int-key")
		value, _ = m.Load(123)
		if value != "int-key" {
			t.Fatalf("Store 不同类型的键后 Load 应返回对应值，期望 %v，实际返回 %v", "int-key", value)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		m := NewMap()

		// 测试删除不存在的键
		m.Delete("不存在的键") // 不应引发错误

		// 测试删除存在的键
		m.Store("key1", "value1")
		m.Store("key2", "value2")
		m.Store("key3", "value3")

		m.Delete("key2")
		if _, loaded := m.Load("key2"); loaded {
			t.Fatal("删除后 Load 应返回 nil")
		}

		// 确保其他键值对不受影响
		value1, _ := m.Load("key1")
		value3, _ := m.Load("key3")
		if value1 != "value1" || value3 != "value3" {
			t.Fatal("删除操作影响了其他键值对")
		}
	})

	t.Run("LoadOrStore", func(t *testing.T) {
		m := NewMap()

		// 测试存储新键值对
		actual, loaded := m.LoadOrStore("key1", "value1")
		if loaded {
			t.Fatal("LoadOrStore 新键应返回 loaded=false")
		}
		if actual != "value1" {
			t.Fatalf("LoadOrStore 新键应返回提供的值，期望 %v，实际返回 %v", "value1", actual)
		}

		// 测试加载现有键值对
		actual, loaded = m.LoadOrStore("key1", "value1-new")
		if !loaded {
			t.Fatal("LoadOrStore 现有键应返回 loaded=true")
		}
		if actual != "value1" {
			t.Fatalf("LoadOrStore 现有键应返回现有值，期望 %v，实际返回 %v", "value1", actual)
		}

		// 确认值未被更新
		value, _ := m.Load("key1")
		if value != "value1" {
			t.Fatalf("LoadOrStore 不应更新现有值，期望 %v，实际为 %v", "value1", value)
		}
	})

	t.Run("LoadAndDelete", func(t *testing.T) {
		m := NewMap()

		// 测试删除不存在的键
		value, loaded := m.LoadAndDelete("不存在的键")
		if loaded {
			t.Fatal("LoadAndDelete 不存在的键应返回 loaded=false")
		}
		if value != nil {
			t.Fatalf("LoadAndDelete 不存在的键应返回 nil，实际返回 %v", value)
		}

		// 测试删除存在的键
		m.Store("key1", "value1")
		value, loaded = m.LoadAndDelete("key1")
		if !loaded {
			t.Fatal("LoadAndDelete 存在的键应返回 loaded=true")
		}
		if value != "value1" {
			t.Fatalf("LoadAndDelete 应返回被删除的值，期望 %v，实际返回 %v", "value1", value)
		}

		// 确认键已被删除
		if _, loaded := m.Load("key1"); loaded {
			t.Fatal("LoadAndDelete 后键应被删除")
		}

		// 测试多个键值对
		m.Store("key1", "value1")
		m.Store("key2", "value2")
		m.Store("key3", "value3")

		value, _ = m.LoadAndDelete("key2")
		if value != "value2" {
			t.Fatalf("LoadAndDelete 应返回正确的值，期望 %v，实际返回 %v", "value2", value)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		m := NewMap()

		// 存储多个键值对
		for i := range 100 {
			m.Store(fmt.Sprintf("key%d", i), i)
		}

		// 清空映射
		m.Clear()

		// 确认所有键都已被删除
		for i := range 100 {
			if _, loaded := m.Load(fmt.Sprintf("key%d", i)); loaded {
				t.Fatalf("Clear 后键 %s 仍存在", fmt.Sprintf("key%d", i))
			}
		}
	})

	t.Run("Range", func(t *testing.T) {
		m := NewMap()

		// 存储多个键值对
		expectedKeys := []string{"key1", "key2", "key3"}
		expectedValues := []string{"value1", "value2", "value3"}

		for i, key := range expectedKeys {
			m.Store(key, expectedValues[i])
		}

		// 使用 Range 遍历
		visitedKeys := make([]any, 0)
		visitedValues := make([]any, 0)

		m.Range(func(key any, value any) bool {
			visitedKeys = append(visitedKeys, key)
			visitedValues = append(visitedValues, value)
			return true
		})

		// 验证是否遍历了所有键值对
		keyMap := make(map[any]bool)
		valueMap := make(map[any]bool)

		for i := range expectedKeys {
			keyMap[expectedKeys[i]] = true
			valueMap[expectedValues[i]] = true
		}

		for i := range visitedKeys {
			if !keyMap[visitedKeys[i]] {
				t.Fatalf("Range 遍历了意外的键: %s", visitedKeys[i])
			}
			if !valueMap[visitedValues[i]] {
				t.Fatalf("Range 遍历了意外的值: %s", visitedValues[i])
			}
		}

		// 测试提前终止遍历
		count := 0
		m.Range(func(key any, value any) bool {
			count++
			return count < 2 // 只遍历两个元素
		})

		if count != 2 {
			t.Fatalf("Range 提前终止遍历失败，期望遍历 2 个元素，实际遍历 %d 个", count)
		}
	})

	t.Run("RangeConcurrent", func(t *testing.T) {
		m := NewMap()

		// 测试空映射
		var emptyCount int
		m.RangeConcurrent(func(chunk int, key any, value any) bool {
			emptyCount++
			return true
		})
		if emptyCount != 0 {
			t.Fatal("RangeConcurrent 空映射应该不调用处理函数")
		}

		// 填充测试数据
		const dataSize = 10000
		for i := range dataSize {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			m.Store(key, value)
		}

		// 测试并发遍历
		var mu sync.Mutex
		visited := make(map[string]string)
		workerCalled := false

		m.RangeConcurrent(func(chunk int, key any, value any) bool {
			k := key.(string)
			v := value.(string)

			// 验证键值对的正确性
			expectedValue := fmt.Sprintf("value-%s", strings.TrimPrefix(k, "key-"))
			if v != expectedValue {
				t.Errorf("RangeConcurrent 遍历的值不正确，键 %s 期望值 %s，实际值 %s", k, expectedValue, v)
			}

			// 记录已访问的键和值
			mu.Lock()
			visited[k] = v
			mu.Unlock()

			return true
		}, func(worker int) {
			workerCalled = true
			if worker != getShardCount() {
				t.Errorf("工作线程数量应该和 getShardCount 函数返回的数量相等，实际为 %d", worker)
			}
		})

		// 验证是否调用了 worker 回调函数
		if !workerCalled {
			t.Error("RangeConcurrent 应该调用 worker 回调函数")
		}

		// 验证是否遍历了所有元素
		if len(visited) != dataSize {
			t.Errorf("RangeConcurrent 应该遍历所有元素，期望 %d 个，实际遍历 %d 个", dataSize, len(visited))
		}

		// 验证所有键值对是否正确
		for i := range dataSize {
			key := fmt.Sprintf("key-%d", i)
			expectedValue := fmt.Sprintf("value-%d", i)
			if value, ok := visited[key]; !ok || value != expectedValue {
				t.Errorf("键 %s 的值不正确，期望 %s，实际 %s，存在状态: %v", key, expectedValue, value, ok)
			}
		}

		// 测试提前终止遍历
		count := 0
		maxCount := dataSize / 2
		m.RangeConcurrent(func(chunk int, key any, value any) bool {
			mu.Lock()
			count++
			shouldContinue := count < maxCount
			mu.Unlock()
			return shouldContinue
		})

		// 由于并发执行，实际访问的元素数可能会略多于 maxCount
		if count < maxCount {
			t.Errorf("RangeConcurrent 提前终止遍历失败，期望至少 %d 个元素，实际遍历 %d 个", maxCount, count)
		}

		// 测试 nil 处理函数
		m.RangeConcurrent(nil)
		// 如果没有 panic，则测试通过
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		m := NewMap()
		const goroutines = 100
		const operationsPerGoroutine = 10000

		var wg sync.WaitGroup
		wg.Add(goroutines * 2) // 读和写各一半

		// 写入 goroutines
		for i := range goroutines {
			go func(id int) {
				defer wg.Done()
				for j := range operationsPerGoroutine {
					key := fmt.Sprintf("key-%d-%d", id, j)
					value := fmt.Sprintf("value-%d-%d", id, j)
					m.Store(key, value)
				}
			}(i)
		}

		// 读取 goroutines
		for i := range goroutines {
			go func(id int) {
				defer wg.Done()
				for j := range operationsPerGoroutine {
					// 随机选择一个操作
					op := rand.Intn(7)
					key := fmt.Sprintf("key-%d-%d", id, j)

					switch op {
					case 0:
						// Load
						m.Load(key)
					case 1:
						// Delete
						m.Delete(key)
					case 2:
						// LoadOrStore
						m.LoadOrStore(key, fmt.Sprintf("new-value-%d-%d", id, j))
					case 3:
						// LoadAndDelete
						m.LoadAndDelete(key)
					case 4:
						// Clear
						m.Clear()
					case 5:
						// Range
						m.Range(func(key any, value any) bool {
							return true
						})
					case 6:
						// RangeConcurrent
						m.RangeConcurrent(func(chunk int, key any, value any) bool {
							return true
						})
					}
				}
			}(i)
		}

		wg.Wait()
		// 如果没有死锁或 panic，则测试通过
	})
}
