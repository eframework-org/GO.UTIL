// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"os"
	"strings"
	"sync"
)

var (
	// argsCache 缓存已解析的命令行参数。
	// 使用懒加载方式初始化，确保只在首次使用时解析。
	argsCache map[string]string

	// argsCacheLock 保护参数缓存的互斥锁。
	// 用于确保并发访问时的数据一致性。
	argsCacheLock sync.RWMutex

	// argsCacheInit 确保参数缓存只初始化一次。
	// 使用 sync.Once 保证线程安全的一次性初始化。
	argsCacheInit sync.Once
)

// GetArg 获取命令行参数值。
// 参数 key 是要查找的参数名。
// 格式：--key=value 或 --key value，例如：--config=dev.json 或 --config dev.json
// 如果参数不存在，返回空字符串。
func GetArg(key string) string {
	// 使用缓存的参数
	args := GetArgs()
	return args[key]
}

// GetArgs 获取所有命令行参数作为键值对映射。
// 格式：--key1=value1 --key2 value2
// 返回参数名到参数值的映射。
// 该函数是线程安全的，确保参数只被解析一次。
func GetArgs() map[string]string {
	// 确保只初始化一次
	argsCacheInit.Do(func() {
		argsCache = make(map[string]string)
		parseArgs()
	})

	// 返回缓存的结果
	argsCacheLock.RLock()
	defer argsCacheLock.RUnlock()
	return argsCache
}

// parseArgs 解析命令行参数并存储在缓存中。
// 支持两种格式：
// 1. --key=value：直接通过等号连接的键值对
// 2. --key value：通过空格分隔的键值对
// 该函数在并发环境下是安全的，通过互斥锁保护缓存的写入。
func parseArgs() {
	argsCacheLock.Lock()
	defer argsCacheLock.Unlock()

	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		key := strings.TrimPrefix(arg, "--")

		// 处理 --key=value 格式
		if idx := strings.Index(key, "="); idx != -1 {
			argsCache[key[:idx]] = key[idx+1:]
			continue
		}

		// 处理 --key value 格式
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			argsCache[key] = args[i+1]
			i++ // 跳过下一个参数，因为它是值
		}
	}
}
