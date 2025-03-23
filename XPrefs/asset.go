// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XPrefs

import (
	"fmt"
	"strings"
)

// 资产首选项的默认路径.
const assetFile = "Assets/Preferences.json"

// 管理只读的资产首选项.
type prefsAsset struct{ prefsBase }

// read 函数从指定的文件中读取偏好设置。
// 如果没有指定文件，则从默认的资产文件中读取。
// 如果文件读取成功，则返回 true，否则返回 false。
func (pa *prefsAsset) read(file ...string) bool {
	var data []byte
	var err error
	filename := assetFile

	if len(file) > 0 && fileExists(file[0]) {
		filename = file[0]
	} else if !fileExists(assetFile) {
		fmt.Printf("XPrefs.Asset.Read: file %s was not found.\n", assetFile)
		return false
	}

	data, err = readFile(filename)
	if err != nil {
		fmt.Printf("XPrefs.Asset.Read: failed to read file %s: %v\n", filename, err)
		return false
	}
	fmt.Printf("XPrefs.Asset.Read: reading %s.\n", filename)

	return pa.parse(data)
}

// parse 函数解析配置数据并应用命令行参数覆盖。
// 输入字节数组形式的配置数据，解析成功返回 true，失败返回 false。
// 解析完成后会检查命令行参数中以 "Prefs@Asset." 开头的配置项，并用其值覆盖相应的配置。
func (pa *prefsAsset) parse(data []byte) bool {
	defer func() {
		args := parseArgs()
		for k, v := range args {
			if strings.HasPrefix(k, "Prefs@Asset.") {
				key := strings.TrimPrefix(k, "Prefs@Asset.")
				if strings.Contains(key, ".") {
					parts := strings.Split(key, ".")
					current := (any)(pa).(IBase)
					for i := 0; i < len(parts)-1; i++ {
						part := parts[i]
						if !current.Has(part) {
							current.Set(part, New())
						}
						current = current.Get(part).(IBase)
					}
					current.Set(parts[len(parts)-1], v)
				} else {
					pa.Set(key, v)
				}
				fmt.Printf("XPrefs.Asset.Parse: override %s = %s\n", key, v)
			}
		}
	}()

	return pa.prefsBase.parse(data)
}
