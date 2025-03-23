// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XFile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eframework-org/GO.UTIL/XString"
)

// Separator 是路径分隔符（POSIX风格）。
const Separator string = "/"

// HasFile 检查文件是否存在。
// 返回 true 表示文件存在，false 表示文件不存在。
func HasFile(file string) bool {
	_, e := os.Stat(file)
	return e == nil
}

// OpenFile 打开文件并返回其内容的字节切片。
// 如果文件不存在或读取失败，返回 nil。
func OpenFile(file string) []byte {
	b, e := os.ReadFile(file)
	if e != nil {
		fmt.Printf("XFile.OpenFile: %v\n", e)
		return nil
	}
	return b
}

// SaveFile 将数据保存到指定文件，并设置给定的权限。
// 如果未指定权限，默认使用 os.ModeAppend。
// 返回写入过程中的错误，如果成功则返回 nil。
func SaveFile(file string, data []byte, perm ...os.FileMode) error {
	_perm := os.ModeAppend
	if len(perm) == 1 {
		_perm = perm[0]
	}
	return os.WriteFile(file, data, _perm)
}

// DeleteFile 删除指定的文件。
// 返回 true 表示删除成功，false 表示删除失败。
func DeleteFile(file string) bool {
	if e := os.Remove(file); e != nil {
		fmt.Printf("XFile.DeleteFile: %v\n", e)
		return false
	}
	return true
}

// OpenText 打开文件并返回其内容的字符串。
// 如果文件不存在或读取失败，返回空字符串。
func OpenText(file string) string {
	b, e := os.ReadFile(file)
	if e != nil {
		fmt.Printf("XFile.OpenText: %v\n", e)
		return ""
	}
	return XString.FromBuffer(b)
}

// SaveText 将文本数据保存到指定文件，并设置给定的权限。
// 如果未指定权限，默认使用 os.ModeAppend。
// 返回写入过程中的错误，如果成功则返回 nil。
func SaveText(file string, data string, perm ...os.FileMode) error {
	_perm := os.ModeAppend
	if len(perm) == 1 {
		_perm = perm[0]
	}
	b := XString.ToBuffer(data)
	return os.WriteFile(file, b, _perm)
}

// HasDirectory 检查文件夹是否存在。
// createIsNotExist 为 true 时，如果目录不存在则创建。
// 返回 true 表示目录存在或创建成功，false 表示目录不存在或创建失败。
func HasDirectory(dir string, createIsNotExist ...bool) bool {
	s, e := os.Stat(dir)
	if e != nil {
		if len(createIsNotExist) == 1 && createIsNotExist[0] {
			return CreateDirectory(dir)
		}
	} else {
		if !s.IsDir() {
			if len(createIsNotExist) == 1 && createIsNotExist[0] {
				DeleteFile(dir)
				err := os.MkdirAll(dir, os.ModePerm)
				if err == nil {
					return true
				}
			}
		} else {
			return true
		}
	}
	return false
}

// CreateDirectory 在指定路径创建一个文件夹。
// 如果文件夹已经存在，则不执行任何操作并返回 true。
// 返回 true 表示创建成功或目录已存在，false 表示创建失败。
func CreateDirectory(path string) bool {
	if HasDirectory(path) {
		return true
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Printf("XFile.CreateDirectory: %v\n", err)
		return false
	}
	return true
}

// DirectoryName 返回指定路径的父目录。
// 返回标准化后的父目录路径。
func DirectoryName(path string) string {
	return NormalizePath(filepath.Dir(path))
}

// PathJoin 连接多个路径片段。
// 返回标准化后的完整路径。
func PathJoin(paths ...string) string {
	path := filepath.Join(paths...)
	path = NormalizePath(path)
	return path
}

// NormalizePath 标准化路径字符串。
// 处理特殊前缀（file://、jar:file://）。
// 统一使用 POSIX 风格的路径分隔符。
// 处理 . 和 .. 路径片段。
// 返回标准化后的路径。
func NormalizePath(path string) string {
	if path == "" {
		return ""
	}

	// 1. 处理特殊前缀
	var prefix string
	switch {
	case strings.HasPrefix(path, "file://"):
		prefix = "file://"
		path = path[7:]
	case strings.HasPrefix(path, "jar:file://"):
		prefix = "jar:file://"
		path = path[11:]
	}

	// 2. 统一分隔符为 POSIX 格式
	path = strings.ReplaceAll(path, "\\", "/")

	// 3. 分割并处理路径
	parts := XString.Split(path, Separator)
	var nparts []string

	for _, part := range parts {
		switch part {
		case "", ".":
			// 如果是根路径，保留空部分
			if len(nparts) == 0 {
				nparts = append(nparts, part)
			}
		case "..":
			if len(nparts) > 0 && nparts[len(nparts)-1] != ".." {
				nparts = nparts[:len(nparts)-1]
			}
		default:
			nparts = append(nparts, part)
		}
	}

	// 4. 重新组合路径并添加前缀
	path = strings.Join(nparts, Separator)
	return prefix + path
}
