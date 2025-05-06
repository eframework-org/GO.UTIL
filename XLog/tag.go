// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLog

import (
	"fmt"
	"strings"
	"sync"

	"github.com/petermattis/goid"
)

// tagPool 是日志标签对象池。
// 用于复用 LogTag 对象，减少内存分配和垃圾回收的压力。
var tagPool = sync.Pool{New: func() any {
	tag := &LogTag{}
	tag.Reset()
	return tag
}}

// tagMap 存储 goroutine 与日志标签的映射关系。
// 使用 goroutine ID 作为键，LogTag 实例作为值，实现每个 goroutine 独立的日志标签上下文。
var tagMap sync.Map

// LogTag 定义日志标签类型。
// 提供了一个线程安全的容器，用于存储和管理日志的元数据信息，如键值对和日志级别。
// 支持动态添加和获取标签信息，可以影响日志的输出行为和格式。
type LogTag struct {
	sync.RWMutex                   // 读写锁，保证并发安全
	level        LevelType         // 日志级别，可覆盖全局日志级别
	key, value   []string          // 标签的键值对列表
	text         string            // 标签的文本表示缓存
	data         map[string]string // 标签的映射表示缓存
	count        int               // 当前存储的键值对数量
	rebuildText  bool              // 标记是否需要重建文本表示
	rebuildData  bool              // 标记是否需要重建映射表示
	pooled       bool              // 标记是否已放入对象池
}

// Reset 重置日志标签的所有字段为初始状态。
// 此方法会清空所有标签信息并重置内部状态标记。
// 通常在将标签对象放回对象池前调用，以确保下次使用时的状态正确。
func (tg *LogTag) Reset() {
	tg.level = LevelUndefined
	tg.text = ""
	tg.data = nil
	tg.count = 0
	tg.rebuildText = true
	tg.rebuildData = true
}

// Set 设置日志标签中的键值对。
// 如果键已存在，则更新其值；如果键不存在，则添加新的键值对。
// 设置后会标记需要重建文本和数据缓存。
func (tg *LogTag) Set(key string, value string) {
	defer tg.Unlock()
	tg.Lock()

	oindex := -1
	if tg.count > 0 {
		for i := 0; i < tg.count; i++ {
			if tg.key[i] == key {
				oindex = i
				break
			}
		}
	}
	if oindex != -1 {
		tg.key[oindex] = key
		tg.value[oindex] = value
	} else {
		if tg.count >= len(tg.key) || len(tg.key) == 0 {
			tg.key = append(tg.key, key)
			tg.value = append(tg.value, value)
		} else {
			tg.key[tg.count] = key
			tg.value[tg.count] = value
		}
		tg.count++
	}
	tg.rebuildText = true
	tg.rebuildData = true
}

// Get 获取日志标签中指定键对应的值，如果键不存在则返回空字符串。
func (tg *LogTag) Get(key string) string {
	defer tg.Unlock()
	tg.Lock()

	if tg.count > 0 {
		for i := 0; i < tg.count; i++ {
			if tg.key[i] == key {
				return tg.value[i]
			}
		}
	}
	return ""
}

// Text 返回日志标签的文本表示。
// 格式为 "[key1 = value1, key2 = value2, ...]"。
// 如果没有标签，则返回空字符串。
// 使用缓存机制，仅在标签内容变更时重建文本。
func (tg *LogTag) Text() string {
	defer tg.Unlock()
	tg.Lock()

	if !tg.rebuildText {
		return tg.text
	} else {
		tg.rebuildText = false
		if tg.count > 0 {
			var builder strings.Builder
			builder.WriteString("[")
			first := true
			for i := range tg.count {
				if !first {
					builder.WriteString(", ")
				} else {
					first = false
				}
				builder.WriteString(fmt.Sprintf("%s = %s", tg.key[i], tg.value[i]))
			}
			builder.WriteString("]")
			tg.text = builder.String()
		} else {
			tg.text = ""
		}
		return tg.text
	}
}

// Data 返回日志标签的映射表示，将所有键值对转换为 map[string]string 格式。
// 使用缓存机制，仅在标签内容变更时重建映射。
func (tg *LogTag) Data() map[string]string {
	defer tg.Unlock()
	tg.Lock()

	if !tg.rebuildData {
		return tg.data
	} else {
		tg.rebuildData = false
		tg.data = make(map[string]string)
		if tg.count > 0 {
			for i := 0; i < tg.count; i++ {
				tg.data[tg.key[i]] = tg.value[i]
			}
		}
		return tg.data
	}
}

// Clone 创建并返回当前日志标签的深度副本。
// 新标签包含原标签的所有键值对，但使用独立的内存空间。
func (tg *LogTag) Clone() *LogTag {
	defer tg.Unlock()
	tg.Lock()

	ntag := GetTag()
	if tg.count > 0 {
		for i := range tg.count {
			ntag.Set(tg.key[i], tg.value[i])
		}
	}
	ntag.level = tg.level

	return ntag
}

// Level 设置或获取日志标签的日志级别。
// 当提供日志级别参数时进行设置，无参数时返回当前级别。
// 标签的日志级别可以覆盖全局的日志级别设置。
func (tg *LogTag) Level(level ...LevelType) LevelType {
	if len(level) > 0 {
		tg.level = level[0]
	}
	return tg.level
}

// GetTag 从对象池中获取一个已重置为初始状态的日志标签实例。
func GetTag() *LogTag {
	tag := tagPool.Get().(*LogTag)
	tag.pooled = false
	return tag
}

// PutTag 将指定的日志标签实例返回到对象池。
// 在返回前会重置标签的状态，并标记其已被放入池中。
// 重复放入同一标签是安全的，会被忽略。
func PutTag(tag *LogTag) {
	if tag != nil && !tag.pooled {
		tag.Reset()
		tag.pooled = true
		tagPool.Put(tag)
	}
}

// Watch 将日志标签与当前 goroutine 关联。
// 如果提供了日志标签参数则使用该标签，否则创建新的标签。
// 此函数实现了 goroutine 本地的日志标签上下文。
func Watch(tag ...*LogTag) *LogTag {
	var tmpTag *LogTag // var tmpTag *LogTag, 临时变量，用于存储最终使用的 LogTag 实例
	if len(tag) > 0 && tag[0] != nil {
		tmpTag = tag[0]
	} else {
		tmpTag = GetTag()
	}
	tagMap.Store(goid.Get(), tmpTag) // 将 LogTag 实例存储到全局 map 中，使用当前 goroutine 的 ID 作为 key

	return tmpTag
}

// Tag 获取或创建与当前 goroutine 关联的日志标签。
// 可以通过可变参数设置键值对，格式为 [key1, value1, key2, value2, ...]。
// 如果当前 goroutine 没有关联的标签且未提供键值对，则返回 nil。
func Tag(pairs ...string) *LogTag {
	if len(pairs) == 0 {
		val, ok := tagMap.Load(goid.Get())
		if ok {
			return val.(*LogTag)
		} else {
			return nil
		}
	} else {
		tmpTag := GetTag()
		val, ok := tagMap.LoadOrStore(goid.Get(), tmpTag)
		if ok {
			PutTag(tmpTag)
		}
		tag := val.(*LogTag)
		kvl := len(pairs)
		if kvl > 0 {
			if kvl == 1 {
				tag.Set(pairs[0], "")
			} else {
				for i := 0; i < kvl; i += 2 {
					tag.Set(pairs[i], pairs[i+1])
				}
			}
		}
		return tag
	}
}

// Defer 清理当前 goroutine 的日志标签。
// 从全局映射中移除当前 goroutine 的标签关联，并将标签实例返回到对象池。
// 通常在 goroutine 结束前调用，用于防止内存泄漏。
func Defer() {
	tag, _ := tagMap.LoadAndDelete(goid.Get())
	if tag != nil {
		PutTag(tag.(*LogTag))
	}
}
