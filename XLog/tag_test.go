// -------------------------------------------------//
//
//	                MIT License                  //
//	Copyright (c) 2025 EFramework Organization   //
//	      SEE LICENSE.md FOR MORE DETAILS        //
//
// -------------------------------------------------//

package XLog

import (
	"fmt"
	"sync"
	"testing"
)

// 测试创建和获取新标签实例。
func TestNewTagAndGetTag(t *testing.T) {
	tag1 := GetTag()
	if tag1 == nil {
		t.Errorf("NewTag returned nil")
	}

	tag2 := GetTag()
	if tag2 == nil {
		t.Errorf("GetTag returned nil")
	}

	if tag1 == tag2 {
		t.Errorf("NewTag and GetTag should return different instances")
	}
}

// 测试在标签中设置和获取值。
func TestTagSetAndGet(t *testing.T) {
	tag := GetTag()
	tag.Set("key1", "value1")
	tag.Set("key2", "value2")

	if tag.Get("key1") != "value1" {
		t.Errorf("Expected value1, got %s", tag.Get("key1"))
	}

	if tag.Get("key2") != "value2" {
		t.Errorf("Expected value2, got %s", tag.Get("key2"))
	}

	if tag.Get("key3") != "" {
		t.Errorf("Expected empty string, got %s", tag.Get("key3"))
	}
}

// 测试标签的字符串表示。
func TestTagStr(t *testing.T) {
	tag := GetTag()
	tag.Set("key1", "value1")
	tag.Set("key2", "value2")

	expected := "[key1=value1, key2=value2]"
	if tag.Text() != expected {
		t.Errorf("Expected %s, got %s", expected, tag.Text())
	}
}

// 测试标签的数据映射。
func TestTagData(t *testing.T) {
	tag := GetTag()
	tag.Set("key1", "value1")
	tag.Set("key2", "value2")

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	data := tag.Data()
	if len(data) != len(expected) {
		t.Errorf("Expected map length %d, got %d", len(expected), len(data))
	}

	for k, v := range expected {
		if data[k] != v {
			t.Errorf("Expected %s=%s, got %s=%s", k, v, k, data[k])
		}
	}
}

// 测试从现有标签创建新标签。
func TestCloneTag(t *testing.T) {
	tag := GetTag()
	tag.Set("key1", "value1")
	tag.Set("key2", "value2")

	cloneTag := tag.Clone()
	if cloneTag == nil {
		t.Errorf("New returned nil")
	}

	if cloneTag.Get("key1") != "value1" {
		t.Errorf("Expected value1, got %s", cloneTag.Get("key1"))
	}

	if cloneTag.Get("key2") != "value2" {
		t.Errorf("Expected value2, got %s", cloneTag.Get("key2"))
	}
}

// 测试在不同场景下的 Watch 和 Tag 函数。
func TestWatchAndTag(t *testing.T) {
	t.Run("SingleThread", func(t *testing.T) {
		tag := GetTag()
		tag.Set("key1", "value1")

		Watch(tag)

		retrievedTag := Tag()
		if retrievedTag == nil {
			t.Errorf("Tag returned nil")
		}

		if retrievedTag.Get("key1") != "value1" {
			t.Errorf("Expected value1, got %s", retrievedTag.Get("key1"))
		}

		Defer()
		retrievedTagAfterDefer := Tag()
		if retrievedTagAfterDefer != nil {
			t.Errorf("Tag should return nil after Defer")
		}
	})

	t.Run("TagWithPairs", func(t *testing.T) {
		tag := Tag("key1", "value1", "key2", "value2")
		if tag == nil {
			t.Errorf("Tag returned nil")
		}

		if tag.Get("key1") != "value1" {
			t.Errorf("Expected value1, got %s", tag.Get("key1"))
		}

		if tag.Get("key2") != "value2" {
			t.Errorf("Expected value2, got %s", tag.Get("key2"))
		}
	})

	t.Run("Concurrent", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10

		for i := range numGoroutines {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				tag := GetTag()
				tag.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))

				Watch(tag)

				retrievedTag := Tag()
				if retrievedTag == nil {
					t.Errorf("Tag returned nil in goroutine %d", id)
				}

				if retrievedTag.Get(fmt.Sprintf("key%d", id)) != fmt.Sprintf("value%d", id) {
					t.Errorf("Expected value%d, got %s in goroutine %d", id, retrievedTag.Get(fmt.Sprintf("key%d", id)), id)
				}

				Defer()
				retrievedTagAfterDefer := Tag()
				if retrievedTagAfterDefer != nil {
					t.Errorf("Tag should return nil after Defer in goroutine %d", id)
				}
			}(i)
		}

		wg.Wait()
	})
}

// 测试 LogTag 的 Set 方法的性能。
//
// Parameters:
//   - b: 基准测试对象。
func BenchmarkTagSet(b *testing.B) {
	tag := GetTag()
	for i := 0; i < b.N; i++ {
		tag.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
}

// 测试 LogTag 的 Get 方法的性能。
//
// Parameters:
//   - b: 基准测试对象。
func BenchmarkTagGet(b *testing.B) {
	tag := GetTag()
	for i := 0; i < b.N; i++ {
		tag.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tag.Get(fmt.Sprintf("key%d", i%b.N))
	}
}

// 测试 LogTag 的 Str 方法的性能。
//
// Parameters:
//   - b: 基准测试对象。
func BenchmarkTagStr(b *testing.B) {
	tag := GetTag()
	for i := 0; i < b.N; i++ {
		tag.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tag.Text()
	}
}

// 测试 LogTag 的 Data 方法的性能。
//
// Parameters:
//   - b: 基准测试对象。
func BenchmarkTagData(b *testing.B) {
	tag := GetTag()
	for i := 0; i < b.N; i++ {
		tag.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tag.Data()
	}
}

// 测试 LogTag 的 Watch 和 Tag 方法在多线程环境下的性能。
//
// Parameters:
//   - b: 基准测试对象。
func BenchmarkWatchAndTag(b *testing.B) {
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go func(id int) {
				defer wg.Done()

				tag := GetTag()
				tag.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))

				Watch(tag).Level(LevelAlert)

				retrievedTag := Tag()
				if retrievedTag == nil {
					b.Errorf("Tag returned nil in goroutine %d", id)
				}

				if retrievedTag.Level() != LevelAlert {
					b.Errorf("Expected LevelAlert, got %d in goroutine %d", retrievedTag.Level(), id)
				}

				if retrievedTag.Get(fmt.Sprintf("key%d", id)) != fmt.Sprintf("value%d", id) {
					b.Errorf("Expected value%d, got %s in goroutine %d", id, retrievedTag.Get(fmt.Sprintf("key%d", id)), id)
				}

				Defer()
				retrievedTagAfterDefer := Tag()
				if retrievedTagAfterDefer != nil {
					b.Errorf("Tag should return nil after Defer in goroutine %d", id)
				}
			}(i)
		}
		wg.Wait()
	}
}
