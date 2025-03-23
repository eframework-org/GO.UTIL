// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XUtility

import (
	"math/rand"
)

// MaxValue 返回两个整数中的最大值。
// a 第一个整数。
// b 第二个整数。
// 返回两个整数中的较大值。
func MaxValue(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

// MinValue 返回两个整数中的最小值。
// a 第一个整数。
// b 第二个整数。
// 返回两个整数中的较小值。
func MinValue(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}

// RandInt 生成指定范围内的随机整数。
// min 最小值（包含）。
// max 最大值（不包含）。
// 返回 [min, max) 范围内的随机整数，如果 min >= max，则返回 max。
func RandInt(min int, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min) + min
}
