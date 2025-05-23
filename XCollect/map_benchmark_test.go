// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XCollect

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

// 基准测试：比较 XCollect.Map 和 sync.Map 的性能
// 多核测试：go test -benchmem -bench=. -cpu=1,2,4,8,16,32
func BenchmarkMap(b *testing.B) {
	b.Run("Mixed", func(b *testing.B) {
		for _, count := range []int{1000, 10000} {
			b.Run(fmt.Sprintf("XCollect.Map/Count-%d", count), func(b *testing.B) {
				m := NewMap()
				// 预填充数据
				for i := range count {
					m.Store(strconv.Itoa(i), i)
				}

				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						op := i % 5 // 5种操作：Load, Store, Delete, LoadOrStore, LoadAndDelete
						key := strconv.Itoa(i % count)

						switch op {
						case 0:
							m.Load(key)
						case 1:
							m.Store(key, i)
						case 2:
							m.Delete(key)
						case 3:
							m.LoadOrStore(key, i)
						case 4:
							m.LoadAndDelete(key)
						}
						i++
					}
				})
			})

			b.Run(fmt.Sprintf("sync.Map/Count-%d", count), func(b *testing.B) {
				var m sync.Map
				// 预填充数据
				for i := range count {
					m.Store(strconv.Itoa(i), i)
				}

				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						op := i % 5 // 5种操作：Load, Store, Delete, LoadOrStore, LoadAndDelete
						key := strconv.Itoa(i % count)

						switch op {
						case 0:
							m.Load(key)
						case 1:
							m.Store(key, i)
						case 2:
							m.Delete(key)
						case 3:
							m.LoadOrStore(key, i)
						case 4:
							m.LoadAndDelete(key)
						}
						i++
					}
				})
			})
		}
	})

	b.Run("Range", func(b *testing.B) {
		for _, count := range []int{1000, 10000, 100000} {
			b.Run(fmt.Sprintf("XCollect.Map/Count-%d", count), func(b *testing.B) {
				m := NewMap()
				for i := range count {
					m.Store(strconv.Itoa(i), i)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					m.Range(func(key, value any) bool {
						return true
					})
				}
			})

			b.Run(fmt.Sprintf("XCollect.Map/Concurrent/Count-%d", count), func(b *testing.B) {
				m := NewMap()
				for i := range count {
					m.Store(strconv.Itoa(i), i)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					m.RangeConcurrent(func(chunk int, key, value any) bool {
						return true
					})
				}
			})

			b.Run(fmt.Sprintf("sync.Map/Count-%d", count), func(b *testing.B) {
				var m sync.Map
				for i := range count {
					m.Store(strconv.Itoa(i), i)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					m.Range(func(key, value any) bool {
						return true
					})
				}
			})

			b.Run(fmt.Sprintf("map/Count-%d", count), func(b *testing.B) {
				m := make(map[string]int)

				for i := range count {
					m[strconv.Itoa(i)] = i
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for range m {
					}
				}
			})
		}
	})
}
