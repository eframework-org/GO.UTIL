// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	// localPath 存储本地数据路径。
	// 默认值为 "Local"，可通过命令行参数 --LocalPath 指定。
	localPath string

	// assetPath 存储资产文件路径。
	// 默认值为 "Assets"，可通过命令行参数 --AssetPath 指定。
	assetPath string

	// pathsOnce 确保路径只初始化一次。
	// 使用 sync.Once 保证线程安全的一次性初始化。
	pathsOnce sync.Once
)

// LocalPath 返回本地数据存储路径。
// 该路径用于存储应用程序的本地数据文件。
// 路径特点：
//   - 自动转换为绝对路径
//   - 统一使用斜杠（/）作为分隔符
//   - 自动创建目录（权限：0755）
//
// 返回值：
//   - string：规范化后的本地数据路径
//
// 如果路径无效或无法创建目录，将会触发 panic。
func LocalPath() string {
	initPaths()
	return localPath
}

// AssetPath 返回资产文件存储路径。
// 该路径用于存储应用程序的资产文件（如图片、配置等）。
// 路径特点：
//   - 自动转换为绝对路径
//   - 统一使用斜杠（/）作为分隔符
//
// 返回值：
//   - string：规范化后的资产文件路径
//
// 如果路径无效，将会触发 panic。
func AssetPath() string {
	initPaths()
	return assetPath
}

// initPaths 初始化本地数据路径和资产文件路径。
// 该函数在首次调用时执行以下操作：
// 1. 从命令行参数获取路径配置
// 2. 转换为绝对路径
// 3. 规范化路径分隔符
// 4. 创建必要的目录
//
// 该函数是线程安全的，确保只执行一次初始化。
// 如果路径无效或无法创建目录，将会触发 panic。
func initPaths() {
	pathsOnce.Do(func() {
		// 初始化本地路径
		localPath = "Local"
		if path := GetArg("LocalPath"); path != "" {
			localPath = path
		}
		if abs, err := filepath.Abs(localPath); err != nil {
			panic(fmt.Sprintf("Invalid local path: %v", err))
		} else {
			localPath = abs
		}
		localPath = filepath.ToSlash(localPath)
		if err := os.MkdirAll(localPath, 0755); err != nil {
			panic(fmt.Sprintf("Failed to create local directory: %v", err))
		}

		// 初始化资产路径
		assetPath = "Assets"
		if path := GetArg("AssetPath"); path != "" {
			assetPath = path
		}
		if abs, err := filepath.Abs(assetPath); err != nil {
			panic(fmt.Sprintf("Invalid asset path: %v", err))
		} else {
			assetPath = abs
		}
		assetPath = filepath.ToSlash(assetPath)
	})
}
