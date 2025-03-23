// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XCollect

// Index 在数组中查找指定元素，返回其第一次出现的索引位置。如果未找到或数组为空，返回 -1。
// 参数 ele 可以是具体值或判断函数 func(T) bool。
func Index[T comparable](arr []T, ele any) int {
	if arr != nil {
		switch v := ele.(type) {
		case T:
			for k, item := range arr {
				if item == v {
					return k
				}
			}
		case func(T) bool:
			for k, item := range arr {
				if v(item) {
					return k
				}
			}
		default:
			return -1
		}
	}
	return -1
}

// Contains 检查指定元素是否存在于数组中。参数 ele 可以是具体值或条件函数 func(T) bool。
// 如果元素存在返回 true，否则返回 false。
func Contains[T comparable](arr []T, ele any) bool {
	if arr != nil {
		return Index(arr, ele) >= 0
	}
	return false
}

// Remove 从数组中移除所有指定的元素。参数 ele 可以是具体值或条件函数 func(T) bool。
// 返回移除元素后的新数组。
func Remove[T comparable](arr []T, ele any) []T {
	if arr != nil {
		for {
			idx := Index(arr, ele)
			if idx >= 0 {
				arr = append(arr[:idx], arr[idx+1:]...)
			} else {
				break
			}
		}
	}
	return arr
}

// Delete 删除数组中指定索引位置的元素。如果索引无效或数组为空，返回原数组。
func Delete[T comparable](arr []T, idx int) []T {
	if arr != nil {
		if idx < len(arr) {
			arr = append(arr[:idx], arr[idx+1:]...)
		}
	}
	return arr
}

// Append 将元素添加到数组末尾，返回添加元素后的新数组。
func Append[T comparable](arr []T, ele T) []T {
	return append(arr, ele)
}

// Insert 在数组的指定位置插入元素。如果索引无效或数组为空，返回原数组。
func Insert[T comparable](arr []T, idx int, ele T) []T {
	if arr == nil || idx < 0 || idx > len(arr) {
		return arr
	}
	arr = append(arr[:idx], append([]T{ele}, arr[idx:]...)...)
	return arr
}
