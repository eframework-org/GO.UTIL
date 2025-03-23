// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XUtility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMaxValue 测试 MaxValue 函数
func TestMaxValue(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"Positive Numbers", 5, 3, 5},
		{"Negative Numbers", -2, -5, -2},
		{"Mixed Numbers", -1, 1, 1},
		{"Equal Numbers", 4, 4, 4},
		{"Zero and Positive", 0, 7, 7},
		{"Zero and Negative", 0, -7, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxValue(tt.a, tt.b)
			assert.Equal(t, tt.expected, result,
				"MaxValue(%d, %d) = %d; want %d",
				tt.a, tt.b, result, tt.expected)
		})
	}
}

// TestMinValue 测试 MinValue 函数
func TestMinValue(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"Positive Numbers", 5, 3, 3},
		{"Negative Numbers", -2, -5, -5},
		{"Mixed Numbers", -1, 1, -1},
		{"Equal Numbers", 4, 4, 4},
		{"Zero and Positive", 0, 7, 0},
		{"Zero and Negative", 0, -7, -7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinValue(tt.a, tt.b)
			assert.Equal(t, tt.expected, result,
				"MinValue(%d, %d) = %d; want %d",
				tt.a, tt.b, result, tt.expected)
		})
	}
}

// TestRandInt 测试 RandInt 函数
func TestRandInt(t *testing.T) {
	t.Run("Normal Range", func(t *testing.T) {
		min, max := 1, 10
		for i := 0; i < 100; i++ {
			result := RandInt(min, max)
			assert.GreaterOrEqual(t, result, min,
				"RandInt(%d, %d) = %d; should be >= %d",
				min, max, result, min)
			assert.Less(t, result, max,
				"RandInt(%d, %d) = %d; should be < %d",
				min, max, result, max)
		}
	})

	t.Run("Zero Range", func(t *testing.T) {
		min, max := 0, 1
		for i := 0; i < 100; i++ {
			result := RandInt(min, max)
			assert.Equal(t, 0, result,
				"RandInt(%d, %d) = %d; should be 0",
				min, max, result)
		}
	})

	t.Run("Single Value", func(t *testing.T) {
		min, max := 5, 6
		for i := 0; i < 100; i++ {
			result := RandInt(min, max)
			assert.Equal(t, 5, result,
				"RandInt(%d, %d) = %d; should be 5",
				min, max, result)
		}
	})

	t.Run("Invalid Range", func(t *testing.T) {
		tests := []struct {
			name string
			min  int
			max  int
		}{
			{"Equal Values", 5, 5},
			{"Reversed Range", 10, 1},
			{"Negative Range", -5, -10},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := RandInt(tt.min, tt.max)
				assert.Equal(t, tt.max, result,
					"RandInt(%d, %d) = %d; should return max value %d for invalid range",
					tt.min, tt.max, result, tt.max)
			})
		}
	})

	t.Run("Distribution", func(t *testing.T) {
		min, max := 0, 10
		counts := make(map[int]int)
		iterations := 10000

		// 生成大量随机数并统计分布
		for i := 0; i < iterations; i++ {
			result := RandInt(min, max)
			counts[result]++
		}

		// 检查是否所有可能的值都出现过
		for i := min; i < max; i++ {
			count := counts[i]
			assert.Greater(t, count, 0,
				"Value %d was never generated in %d iterations",
				i, iterations)
		}

		// 检查是否没有超出范围的值
		for value := range counts {
			assert.GreaterOrEqual(t, value, min,
				"Generated value %d is less than minimum %d",
				value, min)
			assert.Less(t, value, max,
				"Generated value %d is greater than or equal to maximum %d",
				value, max)
		}
	})
}
