// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"

	"github.com/eframework-org/GO.UTIL/XPrefs"
	"github.com/eframework-org/GO.UTIL/XString"
)

const (
	// ModeUnknown 表示未知的运行模式。
	ModeUnknown = "Unknown"
	// ModeDev 表示开发环境。
	ModeDev = "Dev"
	// ModeTest 表示测试环境。
	ModeTest = "Test"
	// ModeStaging 表示预发布环境。
	ModeStaging = "Staging"
	// ModeProd 表示生产环境。
	ModeProd = "Prod"
)

const (
	// AppUnknown 表示未知的应用类型。
	AppUnknown = "Unknown"
	// AppServer 表示服务器应用。
	AppServer = "Server"
	// AppClient 表示客户端应用。
	AppClient = "Client"
)

const (
	// PlatformUnknown 表示未知的运行平台。
	PlatformUnknown = "Unknown"
	// PlatformWindows 表示 Windows 平台。
	PlatformWindows = "Windows"
	// PlatformLinux 表示 Linux 平台。
	PlatformLinux = "Linux"
	// PlatformOSX 表示 macOS 平台。
	PlatformOSX = "OSX"
	// PlatformAndroid 表示 Android 平台。
	PlatformAndroid = "Android"
	// PlatformiOS 表示 iOS 平台。
	PlatformiOS = "iOS"
	// PlatformBrowser 表示浏览器平台。
	PlatformBrowser = "Browser"
)

const (
	// PrefsApp 应用类型的配置键。
	PrefsApp = "Env/App"
	// PrefsAppDefault 应用类型的默认值。
	PrefsAppDefault = AppServer
	// PrefsMode 运行模式的配置键。
	PrefsMode = "Env/Mode"
	// PrefsModeDefault 运行模式的默认值。
	PrefsModeDefault = ModeDev
	// PrefsSolution 解决方案的配置键。
	PrefsSolution = "Env/Solution"
	// PrefsSolutionDefault 解决方案的默认值。
	PrefsSolutionDefault = "Unknown"
	// PrefsProject 项目的配置键。
	PrefsProject = "Env/Project"
	// PrefsProjectDefault 项目的默认值。
	PrefsProjectDefault = "Unknown"
	// PrefsProduct 产品的配置键。
	PrefsProduct = "Env/Product"
	// PrefsProductDefault 产品的默认值。
	PrefsProductDefault = "Unknown"
	// PrefsChannel 渠道的配置键。
	PrefsChannel = "Env/Channel"
	// PrefsChannelDefault 渠道的默认值。
	PrefsChannelDefault = "Default"
	// PrefsVersion 版本的配置键。
	PrefsVersion = "Env/Version"
	// PrefsVersionDefault 版本的默认值。
	PrefsVersionDefault = "0.0.0"
	// PrefsAuthor 作者的配置键。
	PrefsAuthor = "Env/Author"
	// PrefsSecret 密钥的配置键。
	PrefsSecret = "Env/Secret"
	// PrefsRemote 远程配置的配置键。
	PrefsRemote = "Env/Remote"
	// PrefsRemoteDefault 远程配置的默认值。
	PrefsRemoteDefault = "${Env.OssPublic}/Prefs/${Env.Solution}/${Env.Channel}/${Env.Platform}/${Env.Version}/Preferences.json"
)

var (
	// PrefsAuthorDefault 作者的默认值。
	// 默认使用当前系统用户名。
	PrefsAuthorDefault = func() string {
		data, err := user.Current()
		if err != nil {
			return "Unknown"
		}
		username := data.Username
		// Windows: DOMAIN\username
		if i := strings.LastIndex(username, "\\"); i >= 0 {
			return username[i+1:]
		}
		// Linux/OSX: domain/username 或 username
		if i := strings.LastIndex(username, "/"); i >= 0 {
			return username[i+1:]
		}
		return username
	}()

	// PrefsSecretDefault 密钥的默认值。
	// 默认生成 8 位随机字符串。
	PrefsSecretDefault = XString.Random("N")[:8]
)

var (
	// 缓存标志和值
	bApp      = false
	app       = ""
	bMode     = false
	mode      = ""
	bSolution = false
	solution  = ""
	bProject  = false
	project   = ""
	bProduct  = false
	product   = ""
	bChannel  = false
	channel   = ""
	bVersion  = false
	version   = ""
	bAuthor   = false
	author    = ""
	bSecret   = false
	secret    = ""
	bRemote   = false
	remote    = ""
)

// Platform 返回当前运行平台。
// 支持的平台：
//   - Windows：windows
//   - Linux：linux
//   - macOS：darwin
//   - Android：android
//   - iOS：ios
//   - 浏览器：js
//
// 返回值：
//   - string：平台标识符
func Platform() string {
	if runtime.GOOS == "windows" {
		return PlatformWindows
	} else if runtime.GOOS == "linux" {
		return PlatformLinux
	} else if runtime.GOOS == "darwin" {
		return PlatformOSX
	} else if runtime.GOOS == "android" {
		return PlatformAndroid
	} else if runtime.GOOS == "ios" {
		return PlatformiOS
	} else if runtime.GOOS == "js" {
		return PlatformBrowser
	}
	return PlatformUnknown
}

// App 返回应用程序类型。
// 支持的类型：
//   - Server：服务器应用
//   - Client：客户端应用
//
// 返回值：
//   - string：应用程序类型
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
func App() string {
	if !bApp {
		bApp = true
		app = XPrefs.GetString(PrefsApp, PrefsAppDefault)
		switch app {
		case AppServer, AppClient:
		default:
			app = AppUnknown
		}
	}
	return app
}

// Mode 返回运行模式。
// 支持的模式：
//   - Dev：开发环境
//   - Test：测试环境
//   - Staging：预发布环境
//   - Prod：生产环境
//
// 返回值：
//   - string：运行模式
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
func Mode() string {
	if !bMode {
		bMode = true
		mode = XPrefs.GetString(PrefsMode, PrefsModeDefault)
		switch mode {
		case ModeDev, ModeTest, ModeStaging, ModeProd:
		default:
			mode = ModeUnknown
		}
	}
	return mode
}

// Solution 返回解决方案名称。
// 返回值：
//   - string：解决方案名称
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Solution() string {
	if !bSolution {
		bSolution = true
		solution = vars.Eval(XPrefs.GetString(PrefsSolution, PrefsSolutionDefault))
	}
	return solution
}

// Project 返回项目名称。
// 返回值：
//   - string：项目名称
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Project() string {
	if !bProject {
		bProject = true
		project = vars.Eval(XPrefs.GetString(PrefsProject, PrefsProjectDefault))
	}
	return project
}

// Product 返回产品名称。
// 返回值：
//   - string：产品名称
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Product() string {
	if !bProduct {
		bProduct = true
		product = vars.Eval(XPrefs.GetString(PrefsProduct, PrefsProductDefault))
	}
	return product
}

// Channel 返回渠道名称。
// 返回值：
//   - string：渠道名称
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Channel() string {
	if !bChannel {
		bChannel = true
		channel = vars.Eval(XPrefs.GetString(PrefsChannel, PrefsChannelDefault))
	}
	return channel
}

// Version 返回版本号。
// 返回值：
//   - string：版本号
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Version() string {
	if !bVersion {
		bVersion = true
		version = vars.Eval(XPrefs.GetString(PrefsVersion, PrefsVersionDefault))
	}
	return version
}

// Author 返回作者名称。
// 返回值：
//   - string：作者名称
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Author() string {
	if !bAuthor {
		bAuthor = true
		author = vars.Eval(XPrefs.GetString(PrefsAuthor, PrefsAuthorDefault))
	}
	return author
}

// Secret 返回密钥。
// 返回值：
//   - string：密钥
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
// 如果未设置密钥，将使用随机生成的默认值。
func Secret() string {
	if !bSecret {
		bSecret = true
		if !XPrefs.HasKey(PrefsSecret) {
			fmt.Printf("XEnv.Secret: secret is not set, use default value.\n")
		}
		secret = vars.Eval(XPrefs.GetString(PrefsSecret, PrefsSecretDefault))
	}
	return secret
}

// Remote 返回远程配置文件路径。
// 返回值：
//   - string：远程配置文件路径
//
// 该函数使用懒加载方式获取配置，结果会被缓存。
// 支持环境变量引用求值。
func Remote() string {
	if !bRemote {
		remote = vars.Eval(XPrefs.GetString(PrefsRemote, PrefsRemoteDefault))
		bRemote = true
	}
	return remote
}
