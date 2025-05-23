// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XCollect

import (
	"fmt"
	"hash/fnv"
	"runtime"
	"sync"
	"sync/atomic"
)

// Map 提供一个线程安全的键值映射结构。
// 内部使用分段锁实现，将数据分成多个分片，每个分片有独立的锁，减少锁竞争。
// 键值存储采用 map 用于快速索引，切片用于遍历优化。
// 适用于高并发读写场景，大数据遍历性能优于标准库的 sync.Map。
type Map struct {
	mutex     sync.Mutex
	shards    []*mapShard // 存储多个分片
	shardMask uint32      // 分片索引掩码
}

// mapShard 表示 Map 的一个分片。
// 每个分片包含独立的锁和数据存储，用于减少锁竞争。
type mapShard struct {
	mutex  sync.RWMutex // 每个分片独立的读写锁
	keys   map[any]int  // 存储每个 key 在 values 切片中的索引位置
	values []struct {
		key   any
		value any
	} // 以切片形式存储键值对，便于顺序遍历
}

// NewMap 创建并返回一个新的 Map 实例。
// 返回的 Map 实例是未初始化的，将在首次使用时延迟初始化。
func NewMap() *Map {
	m := &Map{}
	return m
}

// getShardCount 返回用于分片的数量。
// 返回的分片数量基于系统可用的逻辑 CPU 数量，并向上调整为最接近的 2 的幂。
// 这样设计可以通过位运算快速定位分片，提高性能。
func getShardCount() int {
	// 获取当前系统可用的逻辑 CPU 数量（可能不是 2 的幂）
	raw := runtime.NumCPU()

	// 将 raw 向上调整为最近的 2 的幂
	n := max(raw, 1)
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++

	return n
}

// initShard 初始化 Map 的分片结构。
// 根据 CPU 数量确定分片数，并创建相应数量的分片。
// 该方法在首次访问 Map 时被调用，实现延迟初始化。
func (m *Map) initShard() {
	shardCount := getShardCount()
	m.shards = make([]*mapShard, shardCount)
	m.shardMask = uint32(shardCount - 1) // 掩码用于位运算
	// 初始化每个分片
	for i := range shardCount {
		m.shards[i] = &mapShard{}
	}
}

// 计算键的哈希值，确定应该存储在哪个分片上。
func (m *Map) getShard(key any) *mapShard {
	switch v := key.(type) {
	case string:
		hash := fnv.New32a()
		hash.Write([]byte(key.(string)))
		return m.shards[hash.Sum32()&m.shardMask]
	case int:
		return m.shards[uint32(v)&m.shardMask]
	case int32:
		return m.shards[uint32(v)&m.shardMask]
	case int64:
		return m.shards[uint32(v)&m.shardMask]
	case uint:
		return m.shards[uint32(v)&m.shardMask]
	case uint32:
		return m.shards[v&m.shardMask]
	case uint64:
		return m.shards[uint32(v)&m.shardMask]
	default:
		hash := fnv.New32a()
		hash.Write(fmt.Appendf(nil, "%v", key))
		return m.shards[hash.Sum32()&m.shardMask]
	}
}

// Load 返回指定 key 对应的值。
// 如果 key 存在，则返回对应的值和 true，否则返回 nil 和 false。
// 该操作是线程安全的，使用读锁保护。
func (m *Map) Load(key any) (value any, ok bool) {
	if m.shards == nil {
		return nil, false
	}

	shard := m.getShard(key)
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()

	if shard.keys == nil {
		return nil, false
	}
	if idx, exists := shard.keys[key]; exists {
		return shard.values[idx].value, true
	}
	return nil, false
}

// Store 设置指定 key 的值。
// 如果 key 已存在，则更新其值；否则插入新的键值对。
// 该操作是线程安全的，使用写锁保护。
func (m *Map) Store(key any, value any) {
	if m.shards == nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		if m.shards == nil {
			m.initShard()
		}
	}

	shard := m.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if shard.keys == nil {
		shard.keys = make(map[any]int, 16)
	}
	if idx, exists := shard.keys[key]; exists {
		shard.values[idx].value = value
	} else {
		shard.keys[key] = len(shard.values)
		shard.values = append(shard.values, struct {
			key   any
			value any
		}{key: key, value: value})
	}
}

// Delete 删除指定 key 及其对应的值。
// 若 key 存在，删除并重排内部切片以保持数据紧凑；否则不做处理。
// 该操作是线程安全的，使用写锁保护。
func (m *Map) Delete(key any) {
	if m.shards == nil {
		return
	}

	shard := m.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if shard.keys == nil {
		return
	}
	if idx, exists := shard.keys[key]; exists {
		delete(shard.keys, key)

		last := len(shard.values) - 1
		if idx < last {
			shard.values[idx] = shard.values[last]
			shard.keys[shard.values[idx].key] = idx
		}
		shard.values = shard.values[:last]
	}
}

// LoadOrStore 返回指定 key 的值，若 key 不存在则写入默认值。
// actual 表示实际存储的值（已存在的或新写入的）。
// loaded 表示 key 已存在，false 表示 key 不存在并已写入新值。
// 该操作是线程安全的，使用写锁保护。
func (m *Map) LoadOrStore(key any, value any) (actual any, loaded bool) {
	if m.shards == nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		if m.shards == nil {
			m.initShard()
		}
	}

	shard := m.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if shard.keys != nil {
		if idx, exists := shard.keys[key]; exists {
			return shard.values[idx].value, true
		}
	} else {
		shard.keys = make(map[any]int, 16)
	}

	shard.keys[key] = len(shard.values)
	shard.values = append(shard.values, struct {
		key   any
		value any
	}{key: key, value: value})
	return value, false
}

// LoadAndDelete 返回指定 key 对应的值并将其从映射中删除。
// value 是 key 对应的值，如果 key 不存在则为 nil。
// loaded 表示 key 是否存在，true 表示 key 存在并已删除，false 表示 key 不存在。
// 该操作是线程安全的，使用写锁保护。
func (m *Map) LoadAndDelete(key any) (value any, loaded bool) {
	if m.shards == nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		if m.shards == nil {
			m.initShard()
		}
	}

	shard := m.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if shard.keys == nil {
		return nil, false
	}
	if idx, exists := shard.keys[key]; exists {
		value = shard.values[idx].value
		delete(shard.keys, key)

		last := len(shard.values) - 1
		if idx < last {
			shard.values[idx] = shard.values[last]
			shard.keys[shard.values[idx].key] = idx
		}
		shard.values = shard.values[:last]

		return value, true
	}
	return nil, false
}

// Clear 清除所有键值对。
// 该方法会重置所有分片的数据，但保留分片结构，适用于需要重用 Map 的场景。
// 该操作是线程安全的，对每个分片使用写锁保护。
func (m *Map) Clear() {
	if m.shards == nil {
		return
	}

	// 清除每个分片
	for _, shard := range m.shards {
		shard.mutex.Lock()
		if shard.keys != nil {
			shard.keys = make(map[any]int, 16)
		}
		shard.values = shard.values[:0]
		shard.mutex.Unlock()
	}
}

// Range 顺序遍历所有键值对，并调用用户提供的 process 函数。
// process 是处理函数，接收 key 和 value 作为参数，返回布尔值，如果 process 返回 false，则提前中断遍历。
// 该操作是线程安全的，对每个分片使用读锁保护。
// 注意：遍历过程中不应修改 Map 内容，否则可能导致不可预期的结果。
func (m *Map) Range(process func(key any, value any) bool) {
	if m.shards == nil || process == nil {
		return
	}

	// 遍历每个分片
	for _, shard := range m.shards {
		shard.mutex.RLock()
		for _, pair := range shard.values {
			if !process(pair.key, pair.value) {
				return
			}
		}
		shard.mutex.RUnlock()
	}
}

// RangeConcurrent 并发遍历所有键值对，内部根据分片数量自动确定协程数量。
// process 是处理函数，接收分片索引、key 和 value 作为参数，返回布尔值，如果 process 返回 false，则通知所有协程终止遍历。
// worker 是可选的回调函数，在启动并发处理前调用，接收分片数量作为参数。
// 该操作是线程安全的，对每个分片使用读锁保护，适用于需要并行处理大量数据的场景，可显著提高处理速度。
// 注意：遍历过程中不应修改 Map 内容，否则可能导致不可预期的结果。
func (m *Map) RangeConcurrent(process func(chunk int, key any, value any) bool, worker ...func(int)) {
	if m.shards == nil || process == nil {
		return
	}

	if len(worker) > 0 && worker[0] != nil {
		worker[0](len(m.shards))
	}

	var wg sync.WaitGroup
	var done int32

	for i, shard := range m.shards {
		wg.Add(1)

		go func(chunk int, shard *mapShard) {
			defer wg.Done()
			shard.mutex.RLock()
			defer shard.mutex.RUnlock()

			for _, pair := range shard.values {
				// 这里直接读 done 变量，不使用原子读，允许一定概率多遍历以降低同步开销
				if done == 1 {
					return
				}
				if !process(chunk, pair.key, pair.value) {
					atomic.StoreInt32(&done, 1)
					return
				}
			}
		}(i, shard)
	}

	wg.Wait()
}
