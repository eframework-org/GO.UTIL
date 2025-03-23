// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XPrefs

import (
	"os"
	"path/filepath"
	"strings"
)

// parseArgs 解析命令行参数。
// 将命令行参数解析为键值对形式的映射。支持两种格式：
// 1. --key=value 格式
// 2. --key value 格式
// 返回解析后的参数映射。
func parseArgs() map[string]string {
	argsMap := make(map[string]string)
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		key := strings.TrimPrefix(arg, "--")
		if idx := strings.Index(key, "="); idx != -1 {
			argsMap[key[:idx]] = key[idx+1:]
			continue
		}
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			argsMap[key] = args[i+1]
			i++ // Skip next argument as it's the value
		}
	}
	return argsMap
}

// fileExists 检查文件是否存在。
// 输入文件路径，如果文件存在且不是目录则返回 true，否则返回 false。
// 用于在读取配置文件前进行检查。
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// readFile 读取文件内容。
// 输入文件路径，返回文件的完整内容和可能的错误。
// 用于读取配置文件。
func readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// writeFile 写入文件内容。
// 输入文件路径和要写入的数据，如果文件所在目录不存在则会创建。
// 写入成功返回 nil，失败返回错误。
func writeFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
