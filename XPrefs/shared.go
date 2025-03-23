// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XPrefs

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/illumitacit/gostd/quit"
)

var (
	// initMu 用于保护初始化过程的互斥锁。
	initMu sync.Mutex
	// initOnce 确保初始化过程只执行一次的同步控制。
	initOnce sync.Once
	// initSig 用于接收系统信号的通道。
	initSig chan os.Signal

	// asset 全局资产配置实例。
	asset *prefsAsset
	// local 全局本地配置实例。
	local *prefsLocal
)

// Asset 获取资产配置实例。
// 返回一个只读的资产配置实例，用于存储应用程序的默认配置和资源文件。
func Asset() *prefsAsset {
	initOnce.Do(setup)
	return asset
}

// Local 获取本地配置实例。
// 返回一个可读写的本地配置实例，用于存储用户的个性化设置，支持动态修改和持久化。
func Local() *prefsLocal {
	initOnce.Do(setup)
	return local
}

// reset 重置初始化状态。
// 仅用于测试目的，重置配置系统的初始化状态。
func reset() {
	initMu.Lock()
	defer initMu.Unlock()

	// 重置同步对象
	initOnce = sync.Once{}

	// 重置信号通道
	if initSig != nil {
		signal.Stop(initSig)
		close(initSig)
		initSig = nil
	}

	// 重置配置实例
	asset = nil
	local = nil
}

// setup 初始化配置系统。
// 读取配置文件，设置信号处理，启用自动保存功能。
func setup() {
	initMu.Lock()
	defer initMu.Unlock()

	args := parseArgs()
	assetFileArg := args["Prefs@Asset"]
	localFileArg := args["Prefs@Local"]

	asset = &prefsAsset{}
	if assetFileArg != "" {
		asset.read(assetFileArg)
	} else {
		asset.read()
	}

	local = &prefsLocal{}
	if localFileArg != "" {
		local.read(localFileArg)
	} else {
		local.read()
	}

	if initSig != nil {
		signal.Stop(initSig)
		close(initSig)
	}
	initSig = make(chan os.Signal, 1)
	signal.Notify(initSig, syscall.SIGTERM, syscall.SIGINT)

	quit.GetWaiter().Add(1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		wg.Done()
		defer func() {
			Local().Save()
			quit.GetWaiter().Done()
		}()
		for {
			select {
			case sig, ok := <-initSig:
				if ok {
					fmt.Printf("XPrefs.Listen: receive signal of %v.\n", sig.String())
				} else {
					fmt.Printf("XPrefs.Listen: channel of signal is closed.\n")
				}
				return
			case <-quit.GetQuitChannel():
				fmt.Println("XPrefs.Listen: receive signal of QUIT.")
				return
			}
		}
	}()
	wg.Wait()
}

// New 创建配置管理实例。
// 返回一个实现了 IBase 接口的新配置管理对象。
func New() IBase { return new(prefsBase) }

// HasKey 检查配置项是否存在。
// 输入键名和可选的配置源列表，如果未提供配置源则在资产配置中查找。
// 返回 true 表示配置项存在，false 表示不存在。
func HasKey(key string, sources ...IBase) bool {
	if len(sources) == 0 {
		return Asset().Has(key)
	}

	for _, source := range sources {
		if source.Has(key) {
			return true
		}
	}

	return false
}

// Get 获取配置项的值。
// 输入键名和可变参数列表，第一个参数为默认值，后续参数为配置源列表。
// 按优先级从配置源中查找值，如果未找到则返回默认值。
func Get(key string, defvalAndSources ...any) any {
	if len(defvalAndSources) == 0 {
		return Asset().Get(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.Get(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().Get(key)
	}

	return defval
}

// Gets 获取配置项的值数组。
// 输入键名和可变参数列表，第一个参数为默认值数组，后续参数为配置源列表。
// 按优先级从配置源中查找值数组，如果未找到则返回默认值数组。
func Gets(key string, defvalAndSources ...any) []any {
	if len(defvalAndSources) == 0 {
		return Asset().Gets(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.Gets(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().Gets(key)
	}

	if val, ok := defval.([]any); ok {
		return val
	}

	return nil
}

// GetInt 获取配置项的整数值。
// 输入键名和可变参数列表，第一个参数为默认值，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为整数，如果未找到或转换失败则返回默认值。
func GetInt(key string, defvalAndSources ...any) int {
	if len(defvalAndSources) == 0 {
		return Asset().GetInt(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetInt(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetInt(key)
	}

	if val, ok := defval.(int); ok {
		return val
	}

	return 0
}

// GetInts 获取配置项的整数数组。
// 输入键名和可变参数列表，第一个参数为默认值数组，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为整数数组，如果未找到或转换失败则返回默认值数组。
func GetInts(key string, defvalAndSources ...any) []int {
	if len(defvalAndSources) == 0 {
		return Asset().GetInts(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetInts(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetInts(key)
	}

	if val, ok := defval.([]int); ok {
		return val
	}

	return nil
}

// GetFloat 获取配置项的浮点数值。
// 输入键名和可变参数列表，第一个参数为默认值，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为浮点数，如果未找到或转换失败则返回默认值。
func GetFloat(key string, defvalAndSources ...any) float32 {
	if len(defvalAndSources) == 0 {
		return Asset().GetFloat(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetFloat(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetFloat(key)
	}

	if val, ok := defval.(float32); ok {
		return val
	}

	return 0
}

// GetFloats 获取配置项的浮点数数组。
// 输入键名和可变参数列表，第一个参数为默认值数组，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为浮点数数组，如果未找到或转换失败则返回默认值数组。
func GetFloats(key string, defvalAndSources ...any) []float32 {
	if len(defvalAndSources) == 0 {
		return Asset().GetFloats(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetFloats(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetFloats(key)
	}

	if val, ok := defval.([]float32); ok {
		return val
	}

	return nil
}

// GetBool 获取配置项的布尔值。
// 输入键名和可变参数列表，第一个参数为默认值，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为布尔值，如果未找到或转换失败则返回默认值。
func GetBool(key string, defvalAndSources ...any) bool {
	if len(defvalAndSources) == 0 {
		return Asset().GetBool(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetBool(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetBool(key)
	}

	if val, ok := defval.(bool); ok {
		return val
	}
	return false
}

// GetBools 获取配置项的布尔值数组。
// 输入键名和可变参数列表，第一个参数为默认值数组，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为布尔值数组，如果未找到或转换失败则返回默认值数组。
func GetBools(key string, defvalAndSources ...any) []bool {
	if len(defvalAndSources) == 0 {
		return Asset().GetBools(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetBools(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetBools(key)
	}

	if val, ok := defval.([]bool); ok {
		return val
	}

	return nil
}

// GetString 获取配置项的字符串值。
// 输入键名和可变参数列表，第一个参数为默认值，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为字符串，如果未找到或转换失败则返回默认值。
func GetString(key string, defvalAndSources ...any) string {
	if len(defvalAndSources) == 0 {
		return Asset().GetString(key)
	}

	defval := defvalAndSources[0]
	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetString(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetString(key)
	}

	if val, ok := defval.(string); ok {
		return val
	}

	return ""
}

// GetStrings 获取配置项的字符串数组。
// 输入键名和可变参数列表，第一个参数为默认值数组，后续参数为配置源列表。
// 按优先级从配置源中查找并转换为字符串数组，如果未找到或转换失败则返回默认值数组。
func GetStrings(key string, defvalAndSources ...any) []string {
	if len(defvalAndSources) == 0 {
		return Asset().GetStrings(key)
	}

	defval := defvalAndSources[0]

	for _, source := range defvalAndSources[1:] {
		if source, ok := source.(IBase); ok {
			if source.Has(key) {
				return source.GetStrings(key)
			}
		}
	}

	if Asset().Has(key) {
		return Asset().GetStrings(key)
	}

	if val, ok := defval.([]string); ok {
		return val
	}

	return nil
}
