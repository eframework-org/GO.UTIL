// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	// 保存原始参数
	originalArgs := os.Args

	// 每个测试用例后重置状态的辅助函数
	resetState := func() {
		// 重置路径相关状态
		localPath = ""
		assetPath = ""
		pathsOnce = sync.Once{}

		// 重置 args 缓存
		argsCacheLock.Lock()
		argsCache = nil
		argsCacheInit = sync.Once{}
		argsCacheLock.Unlock()

		// 恢复原始参数
		os.Args = originalArgs
	}

	// 测试结束后恢复原始状态
	defer resetState()

	t.Run("Default Paths", func(t *testing.T) {
		resetState()

		// 获取默认路径
		local := LocalPath()
		asset := AssetPath()

		// 验证默认值
		absLocal, _ := filepath.Abs("Local")
		absAsset, _ := filepath.Abs("Assets")
		assert.Equal(t, filepath.ToSlash(absLocal), local)
		assert.Equal(t, filepath.ToSlash(absAsset), asset)

		// 验证本地目录已创建
		_, err := os.Stat(local)
		assert.NoError(t, err)
	})

	t.Run("Custom Paths", func(t *testing.T) {
		resetState()

		// 设置自定义路径
		tmpDir := t.TempDir()
		customLocal := filepath.Join(tmpDir, "CustomLocal")
		customAsset := filepath.Join(tmpDir, "CustomAsset")

		os.Args = []string{"test",
			"--LocalPath=" + customLocal,
			"--AssetPath=" + customAsset,
		}

		// 获取路径
		local := LocalPath()
		asset := AssetPath()

		// 验证自定义值
		absLocal, _ := filepath.Abs(customLocal)
		absAsset, _ := filepath.Abs(customAsset)
		assert.Equal(t, filepath.ToSlash(absLocal), local)
		assert.Equal(t, filepath.ToSlash(absAsset), asset)

		// 验证本地目录已创建
		_, err := os.Stat(local)
		assert.NoError(t, err)
	})

	t.Run("Path Caching", func(t *testing.T) {
		resetState()

		// 首次获取
		local1 := LocalPath()
		asset1 := AssetPath()

		// 修改参数
		os.Args = []string{"test",
			"--LocalPath=NewLocal",
			"--AssetPath=NewAsset",
		}

		// 再次获取
		local2 := LocalPath()
		asset2 := AssetPath()

		// 验证缓存生效
		assert.Equal(t, local1, local2)
		assert.Equal(t, asset1, asset2)
	})

	t.Run("Invalid Local Path", func(t *testing.T) {
		resetState()

		// 设置无效路径
		invalidPath := string([]byte{0}) // 使用 NUL 字符创建无效路径
		os.Args = []string{"test", "--LocalPath=" + invalidPath}

		// 验证 panic
		assert.Panics(t, func() {
			LocalPath()
		})
	})
}
