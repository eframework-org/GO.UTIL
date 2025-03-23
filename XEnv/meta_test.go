// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"sync"
	"testing"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/stretchr/testify/assert"
)

func TestMetaValues(t *testing.T) {
	// 重置缓存状态
	bProject = false
	bProduct = false
	bMode = false
	bVersion = false
	bChannel = false
	bApp = false
	bSecret = false
	bRemote = false
	bSolution = false
	bAuthor = false
	project = ""
	product = ""
	mode = ""
	version = ""
	channel = ""
	app = ""
	secret = ""
	remote = ""
	solution = ""
	author = ""

	tests := []struct {
		name     string
		setup    func()
		validate func(t *testing.T)
	}{
		{
			name: "Default Values",
			setup: func() {
				XPrefs.Asset().Unset(PrefsAuthor)
				XPrefs.Asset().Unset(PrefsSecret)
				bAuthor = false
				bSecret = false
				author = ""
				secret = ""
			},
			validate: func(t *testing.T) {
				assert.Equal(t, PrefsAuthorDefault, Author())
				assert.Equal(t, PrefsSecretDefault, Secret())
			},
		},
		{
			name: "Cache Consistency",
			setup: func() {
				// 设置初始值
				XPrefs.Asset().Set(PrefsAuthor, "InitialAuthor")
				XPrefs.Asset().Set(PrefsSecret, "InitialSecret")
				bAuthor = false
				bSecret = false
				author = ""
				secret = ""

				// 获取值以触发缓存
				_ = Author()
				_ = Secret()

				// 更改值
				XPrefs.Asset().Set(PrefsAuthor, "UpdatedAuthor")
				XPrefs.Asset().Set(PrefsSecret, "UpdatedSecret")
			},
			validate: func(t *testing.T) {
				// 验证返回缓存的值而不是更新后的值
				assert.Equal(t, "InitialAuthor", Author())
				assert.Equal(t, "InitialSecret", Secret())
			},
		},
		{
			name: "Custom Values",
			setup: func() {
				XPrefs.Asset().Set(PrefsProject, "TestProject")
				XPrefs.Asset().Set(PrefsProduct, "TestProduct")
				XPrefs.Asset().Set(PrefsMode, ModeTest)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, "TestProject", Project())
				assert.Equal(t, "TestProduct", Product())
				assert.Equal(t, ModeTest, Mode())
			},
		},
		{
			name: "Cache Consistency",
			setup: func() {
				// 首先设置初始值
				XPrefs.Asset().Set(PrefsProject, "InitialProject")
				XPrefs.Asset().Set(PrefsProduct, "InitialProduct")
				XPrefs.Asset().Set(PrefsMode, ModeDev)

				// 获取一次值以初始化缓存
				Project()
				Product()
				Mode()

				// 然后修改值
				XPrefs.Asset().Set(PrefsProject, "UpdatedProject")
				XPrefs.Asset().Set(PrefsProduct, "UpdatedProduct")
				XPrefs.Asset().Set(PrefsMode, ModeProd)
			},
			validate: func(t *testing.T) {
				// 验证返回缓存的值而不是更新后的值
				assert.Equal(t, "InitialProject", Project())
				assert.Equal(t, "InitialProduct", Product())
				assert.Equal(t, ModeDev, Mode())
			},
		},
		{
			name: "Mode Constants",
			setup: func() {
				// 测试所有模式常量
				modes := []string{ModeDev, ModeTest, ModeStaging, ModeProd}
				for _, m := range modes {
					bMode = false
					mode = ""
					XPrefs.Asset().Set(PrefsMode, m)
					assert.Equal(t, m, Mode())
				}

			},
			validate: func(t *testing.T) {
				// 验证模式常量定义
				assert.Equal(t, "Dev", ModeDev)
				assert.Equal(t, "Test", ModeTest)
				assert.Equal(t, "Staging", ModeStaging)
				assert.Equal(t, "Prod", ModeProd)
			},
		},
		{
			name: "Version Values",
			setup: func() {
				XPrefs.Asset().Set(PrefsVersion, "1.2.3")
			},
			validate: func(t *testing.T) {
				assert.Equal(t, "1.2.3", Version())
			},
		},
		{
			name: "Default Version",
			setup: func() {
				// 清除之前的设置
				XPrefs.Asset().Unset(PrefsVersion)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, PrefsVersionDefault, Version())
			},
		},
		{
			name: "Channel Values",
			setup: func() {
				XPrefs.Asset().Set(PrefsChannel, "Beta")
			},
			validate: func(t *testing.T) {
				assert.Equal(t, "Beta", Channel())
			},
		},
		{
			name: "Default Channel",
			setup: func() {
				// 清除之前的设置
				XPrefs.Asset().Unset(PrefsChannel)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, PrefsChannelDefault, Channel())
			},
		},
		{
			name: "Channel Cache Consistency",
			setup: func() {
				// 设置初始值
				XPrefs.Asset().Set(PrefsChannel, "Alpha")

				// 初始化缓存
				Channel()

				// 修改值
				XPrefs.Asset().Set(PrefsChannel, "Beta")
			},
			validate: func(t *testing.T) {
				// 应该返回缓存的值
				assert.Equal(t, "Alpha", Channel())
			},
		},
		{
			name: "App Values",
			setup: func() {
				// 重置状态
				bApp = false
				app = ""
				XPrefs.Asset().Set(PrefsApp, AppClient)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, AppClient, App())
			},
		},
		{
			name: "Default App",
			setup: func() {
				bApp = false
				app = ""
				XPrefs.Asset().Unset(PrefsApp)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, PrefsAppDefault, App())
			},
		},
		{
			name: "Remote Values",
			setup: func() {
				bRemote = false
				remote = ""
				XPrefs.Asset().Set(PrefsRemote, "https://example.com/prefs")
			},
			validate: func(t *testing.T) {
				assert.Equal(t, "https://example.com/prefs", Remote())
			},
		},
		{
			name: "Default Remote",
			setup: func() {
				bRemote = false
				remote = ""
				XPrefs.Asset().Unset(PrefsRemote)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, vars.Eval(PrefsRemoteDefault), Remote())
			},
		},
		{
			name: "Cache Consistency for New Fields",
			setup: func() {
				// 设置初始值
				bApp = false
				bSecret = false
				bRemote = false
				app = ""
				secret = ""
				remote = ""

				// 使用有效的枚举值
				XPrefs.Asset().Set(PrefsApp, AppServer) // 改用有效的枚举值
				XPrefs.Asset().Set(PrefsSecret, "InitialSecret")
				XPrefs.Asset().Set(PrefsRemote, "InitialRemote")

				// 初始化缓存
				App()
				Secret()
				Remote()

				// 修改值
				XPrefs.Asset().Set(PrefsApp, AppClient) // 改用另一个有效的枚举值
				XPrefs.Asset().Set(PrefsSecret, "UpdatedSecret")
				XPrefs.Asset().Set(PrefsRemote, "UpdatedRemote")
			},
			validate: func(t *testing.T) {
				// 验证返回缓存的值而不是更新后的值
				assert.Equal(t, AppServer, App()) // 期望的是初始设置的有效枚举值
				assert.Equal(t, "InitialSecret", Secret())
				assert.Equal(t, "InitialRemote", Remote())
			},
		},
		{
			name: "Solution Values",
			setup: func() {
				bSolution = false
				solution = ""
				XPrefs.Asset().Set(PrefsSolution, "TestSolution")
			},
			validate: func(t *testing.T) {
				assert.Equal(t, "TestSolution", Solution())
			},
		},
		{
			name: "Default Solution",
			setup: func() {
				bSolution = false
				solution = ""
				XPrefs.Asset().Unset(PrefsSolution)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, PrefsSolutionDefault, Solution())
			},
		},
		{
			name: "Solution Cache Consistency",
			setup: func() {
				// 设置初始值
				bSolution = false
				solution = ""
				XPrefs.Asset().Set(PrefsSolution, "InitialSolution")

				// 初始化缓存
				Solution()

				// 修改值
				XPrefs.Asset().Set(PrefsSolution, "UpdatedSolution")
			},
			validate: func(t *testing.T) {
				// 验证返回缓存的值而不是更新后的值
				assert.Equal(t, "InitialSolution", Solution())
			},
		},
		{
			name: "Author Values",
			setup: func() {
				bAuthor = false
				author = ""
				XPrefs.Asset().Set(PrefsAuthor, "TestAuthor")
			},
			validate: func(t *testing.T) {
				assert.Equal(t, "TestAuthor", Author())
			},
		},
		{
			name: "Default Author",
			setup: func() {
				bAuthor = false
				author = ""
				XPrefs.Asset().Unset(PrefsAuthor)
			},
			validate: func(t *testing.T) {
				assert.Equal(t, vars.Eval(PrefsAuthorDefault), Author())
			},
		},
		{
			name: "Author Cache Consistency",
			setup: func() {
				// 设置初始值
				bAuthor = false
				author = ""
				XPrefs.Asset().Set(PrefsAuthor, "InitialAuthor")

				// 初始化缓存
				Author()

				// 修改值
				XPrefs.Asset().Set(PrefsAuthor, "UpdatedAuthor")
			},
			validate: func(t *testing.T) {
				// 验证返回缓存的值而不是更新后的值
				assert.Equal(t, "InitialAuthor", Author())
			},
		},
		{
			name: "App Enum Validation",
			setup: func() {
				tests := []struct {
					input    string
					expected string
				}{
					{AppServer, AppServer},
					{AppClient, AppClient},
					{"InvalidApp", AppUnknown},
					{"server", AppUnknown}, // 大小写敏感
					{"CLIENT", AppUnknown},
					{"", AppUnknown},
				}

				for _, test := range tests {
					bApp = false
					app = ""
					XPrefs.Asset().Set(PrefsApp, test.input)
					assert.Equal(t, test.expected, App(),
						"App value '%s' should return '%s'", test.input, test.expected)
				}
			},
			validate: func(t *testing.T) {
				// 验证常量定义
				assert.Equal(t, "Unknown", AppUnknown)
				assert.Equal(t, "Server", AppServer)
				assert.Equal(t, "Client", AppClient)
			},
		},
		{
			name: "Mode Enum Validation",
			setup: func() {
				tests := []struct {
					input    string
					expected string
				}{
					{ModeDev, ModeDev},
					{ModeTest, ModeTest},
					{ModeStaging, ModeStaging},
					{ModeProd, ModeProd},
					{"InvalidMode", ModeUnknown},
					{"dev", ModeUnknown}, // 大小写敏感
					{"PROD", ModeUnknown},
					{"", ModeUnknown},
				}

				for _, test := range tests {
					bMode = false
					mode = ""
					XPrefs.Asset().Set(PrefsMode, test.input)
					assert.Equal(t, test.expected, Mode(),
						"Mode value '%s' should return '%s'", test.input, test.expected)
				}
			},
			validate: func(t *testing.T) {
				// 验证常量定义
				assert.Equal(t, "Unknown", ModeUnknown)
				assert.Equal(t, "Dev", ModeDev)
				assert.Equal(t, "Test", ModeTest)
				assert.Equal(t, "Staging", ModeStaging)
				assert.Equal(t, "Prod", ModeProd)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置状态
			bProject = false
			bProduct = false
			bMode = false
			bVersion = false
			bChannel = false
			bApp = false
			bSecret = false
			bRemote = false
			bSolution = false
			bAuthor = false
			project = ""
			product = ""
			mode = ""
			version = ""
			channel = ""
			app = ""
			secret = ""
			remote = ""
			solution = ""
			author = ""

			// 运行测试设置
			tt.setup()

			// 运行验证
			tt.validate(t)
		})
	}
}

func TestMetaConcurrency(t *testing.T) {
	// 重置状态
	bProject = false
	bProduct = false
	bMode = false
	bVersion = false
	bChannel = false
	bApp = false
	bSecret = false
	bRemote = false
	bSolution = false
	bAuthor = false
	project = ""
	product = ""
	mode = ""
	version = ""
	channel = ""
	app = ""
	secret = ""
	remote = ""
	solution = ""
	author = ""

	// 设置初始值
	XPrefs.Asset().Set(PrefsProject, "ConcurrentProject")
	XPrefs.Asset().Set(PrefsProduct, "ConcurrentProduct")
	XPrefs.Asset().Set(PrefsMode, ModeDev)
	XPrefs.Asset().Set(PrefsVersion, "2.0.0")
	XPrefs.Asset().Set(PrefsChannel, "Stable")
	XPrefs.Asset().Set(PrefsApp, AppClient)
	XPrefs.Asset().Set(PrefsSecret, "test_secret_key")
	XPrefs.Asset().Set(PrefsRemote, "https://example.com/prefs")
	XPrefs.Asset().Set(PrefsSolution, "TestSolution")
	XPrefs.Asset().Set(PrefsAuthor, "TestAuthor")

	// 预先初始化所有值，避免并发初始化问题
	Solution()
	Project()
	Product()
	Mode()
	Version()
	Channel()
	App()
	Secret()
	Remote()
	Author()

	// 并发测试
	t.Run("Concurrent Access", func(t *testing.T) {
		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				// 并发访问所有方法
				s := Solution()
				p1 := Project()
				p2 := Product()
				m := Mode()
				v := Version()
				c := Channel()
				a := App()
				s2 := Secret()
				r := Remote()
				au := Author()

				// 验证值一致性
				assert.Equal(t, "TestSolution", s)
				assert.Equal(t, "ConcurrentProject", p1)
				assert.Equal(t, "ConcurrentProduct", p2)
				assert.Equal(t, ModeDev, m)
				assert.Equal(t, "2.0.0", v)
				assert.Equal(t, "Stable", c)
				assert.Equal(t, AppClient, a)
				assert.Equal(t, "test_secret_key", s2)
				assert.Equal(t, "https://example.com/prefs", r)
				assert.Equal(t, "TestAuthor", au)
			}()
		}

		wg.Wait()
	})
}
