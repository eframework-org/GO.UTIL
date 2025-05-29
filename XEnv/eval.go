// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/eframework-org/GO.UTIL/XString"
)

var (
	// vars 全局环境变量求值器实例。
	// 使用懒加载方式初始化，确保只创建一个实例。
	vars *envEval

	// varsOnce 确保求值器只初始化一次。
	// 使用 sync.Once 保证线程安全的一次性初始化。
	varsOnce sync.Once
)

// Vars 返回环境变量求值器。
// 该函数是线程安全的，确保全局只有一个求值器实例。
// 返回值：
//   - *envEval：环境变量求值器实例
func Vars() *envEval {
	varsOnce.Do(func() {
		vars = &envEval{}
	})
	return vars
}

// envEval 是环境变量求值器的实现结构体。
// 提供对环境变量表达式的求值功能。
type envEval struct{}

// Eval 对包含环境变量引用的表达式进行求值。
// 支持的格式：${Env.Key}
// 特殊处理：
//   - 嵌套变量：返回 (Nested) 后缀
//   - 循环引用：返回 (Recursive) 后缀
//   - 未知变量：返回 (Unknown) 后缀
//
// 参数：
//   - input：包含环境变量引用的表达式
//
// 返回值：
//   - string：求值后的结果
//
// 示例：
//
//	${Env.APP_ENV} -> "production"
func (ev *envEval) Eval(input string) string {
	pattern := regexp.MustCompile(`\$\{Env\.([^}]+?)\}`)
	visited := make(map[string]bool)

	var replaceFunc func(string) string
	replaceFunc = func(match string) string {
		key := pattern.FindStringSubmatch(match)[1]

		// 1. 检查嵌套变量
		if strings.Contains(key, "${") {
			return fmt.Sprintf("%v(Nested)", match)
		}

		// 2. 检查循环引用
		if visited[key] {
			return fmt.Sprintf("${Env.%v}(Recursive)", key)
		}
		visited[key] = true
		defer delete(visited, key) // 使用 defer 确保清理访问标记

		// 3. 获取变量值
		var value string
		if key == "LocalPath" {
			value = LocalPath()
		} else if key == "AssetPath" {
			value = AssetPath()
		} else if key == "UserName" {
			value = PrefsAuthorDefault
		} else if key == "Platform" {
			value = Platform()
		} else if key == "App" {
			value = App()
		} else if key == "Mode" {
			value = Mode()
		} else if key == "Solution" {
			value = Solution()
		} else if key == "Project" {
			value = Project()
		} else if key == "Product" {
			value = Product()
		} else if key == "Channel" {
			value = Channel()
		} else if key == "Version" {
			value = Version()
		} else if key == "Author" {
			value = Author()
		} else if key == "Secret" {
			value = Secret()
		} else if key == "NumCPU" {
			value = XString.ToString(runtime.NumCPU())
		} else {
			value = GetArg(key)
		}

		if value == "" {
			value = os.Getenv(key)
		}
		if value != "" {
			return pattern.ReplaceAllStringFunc(value, replaceFunc)
		}

		// 4. 处理未知变量
		return fmt.Sprintf("${Env.%v}(Unknown)", key)
	}

	return pattern.ReplaceAllStringFunc(input, replaceFunc)
}
