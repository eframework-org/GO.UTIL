// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XLoom

import (
	"github.com/eframework-org/GO.UTIL/XLog"
)

// RunAsync 异步执行指定的函数。
// callback 为要异步执行的函数。
// recover 为可选的异常恢复标志，如果为 true，则在发生 panic 时会自动重试。
func RunAsync(callback func(), recover ...bool) {
	if callback == nil {
		return
	}
	go func() {
		defer XLog.Caught(false, func(s string, i int) {
			if len(recover) == 1 && recover[0] {
				RunAsync(callback, recover...)
			}
		})
		callback()
	}()
}

// RunAsyncT1 异步执行带一个参数的函数。
// callback 为要异步执行的函数。
// arg1 为函数的参数。
// recover 为可选的异常恢复标志，如果为 true，则在发生 panic 时会自动重试。
func RunAsyncT1[T1 any](callback func(T1), arg1 T1, recover ...bool) {
	if callback == nil {
		return
	}
	go func() {
		defer XLog.Caught(false, func(s string, i int) {
			if len(recover) == 1 && recover[0] {
				RunAsyncT1(callback, arg1, recover...)
			}
		})
		callback(arg1)
	}()
}

// RunAsyncT2 异步执行带两个参数的函数。
// callback 为要异步执行的函数。
// arg1 为第一个参数。
// arg2 为第二个参数。
// recover 为可选的异常恢复标志，如果为 true，则在发生 panic 时会自动重试。
func RunAsyncT2[T1, T2 any](callback func(T1, T2), arg1 T1, arg2 T2, recover ...bool) {
	if callback == nil {
		return
	}
	go func() {
		defer XLog.Caught(false, func(s string, i int) {
			if len(recover) == 1 && recover[0] {
				RunAsyncT2(callback, arg1, arg2, recover...)
			}
		})
		callback(arg1, arg2)
	}()
}

// RunAsyncT3 异步执行带三个参数的函数。
// callback 为要异步执行的函数。
// arg1 为第一个参数。
// arg2 为第二个参数。
// arg3 为第三个参数。
// recover 为可选的异常恢复标志，如果为 true，则在发生 panic 时会自动重试。
func RunAsyncT3[T1, T2, T3 any](callback func(T1, T2, T3), arg1 T1, arg2 T2, arg3 T3, recover ...bool) {
	if callback == nil {
		return
	}
	go func() {
		defer XLog.Caught(false, func(s string, i int) {
			if len(recover) == 1 && recover[0] {
				RunAsyncT3(callback, arg1, arg2, arg3, recover...)
			}
		})
		callback(arg1, arg2, arg3)
	}()
}
