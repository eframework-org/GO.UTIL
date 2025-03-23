// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XPrefs

import (
	"fmt"
	"strings"
)

// localFile 本地首选项的默认路径.
const localFile = "Local/Preferences.json"

// prefsLocal 管理可读写的本地首选项.
type prefsLocal struct {
	prefsBase
	file string // 本地首选项文件路径.
}

// read 函数从指定的文件中读取偏好设置。
// 如果没有指定文件，则从默认的本地偏好设置文件中读取。
// 如果文件读取成功，则返回 true，否则返回 false。
func (pl *prefsLocal) read(file ...string) bool {
	if len(file) > 0 {
		pl.file = file[0]
	} else {
		pl.file = localFile
	}

	fmt.Printf("XPrefs.Local.Read: reading %s.\n", pl.file)

	if !fileExists(pl.file) {
		return false
	}

	data, err := readFile(pl.file)
	if err != nil {
		fmt.Printf("XPrefs.Local.Read: failed to read file %s: %v\n", pl.file, err)
		return false
	}

	return pl.parse(data)
}

// parse 函数解析配置数据并应用命令行参数覆盖。
// 输入字节数组形式的配置数据，解析成功返回 true，失败返回 false。
// 解析完成后会检查命令行参数中以 "Prefs@Local." 开头的配置项，并用其值覆盖相应的配置。
func (pl *prefsLocal) parse(data []byte) bool {
	defer func() {
		args := parseArgs()
		for k, v := range args {
			if strings.HasPrefix(k, "Prefs@Local.") {
				key := strings.TrimPrefix(k, "Prefs@Local.")
				if strings.Contains(key, ".") {
					parts := strings.Split(key, ".")
					current := (any)(pl).(IBase)
					for i := 0; i < len(parts)-1; i++ {
						part := parts[i]
						if !current.Has(part) {
							current.Set(part, New())
						}
						current = current.Get(part).(IBase)
					}
					current.Set(parts[len(parts)-1], v)
				} else {
					pl.Set(key, v)
				}
				fmt.Printf("XPrefs.Local.Parse: override %s = %s\n", key, v)
			}
		}
	}()

	return pl.prefsBase.parse(data)
}

// Save 函数将偏好设置保存到指定的文件中。
// 如果没有指定文件，则保存到默认的本地偏好设置文件中。
// 如果文件保存成功，则返回 true，否则返回 false。
func (pl *prefsLocal) Save(file ...string) bool {
	sfile := pl.file
	if len(file) > 0 {
		sfile = file[0]
	}
	if sfile == "" {
		sfile = localFile
	}

	json := pl.Json(true)
	err := writeFile(sfile, []byte(json))
	if err != nil {
		fmt.Printf("XPrefs.Local.Save: save file err: %v\n", err)
		return false
	}
	fmt.Printf("XPrefs.Local.Save: persisted to %s.\n", sfile)
	return true
}
