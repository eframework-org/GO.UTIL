// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XFile 提供了一组实用的文件系统操作工具函数，支持文件读写和路径处理等功能。

功能特性

  - 文件操作：支持文件的常用操作，如：检查、读取、写入、删除等
  - 目录操作：支持目录的常用操作，如：检查、创建、获取目录名称等
  - 路径处理：统一的 POSIX 风格路径处理，支持路径拼接及标准化

使用手册

1. 文件操作

1.1 文件检查

	exists := XFile.HasFile("path/to/file.txt")

1.2 文件读写

	data := XFile.OpenFile("config.dat")
	err := XFile.SaveFile("config.dat", data, 0644)
	text := XFile.OpenText("config.txt")
	err := XFile.SaveText("config.txt", "Hello World", 0644)

2. 目录操作

2.1 目录管理

	exists := XFile.HasDirectory("path/to/dir")
	success := XFile.CreateDirectory("path/to/dir")
	parent := XFile.DirectoryName("path/to/file.txt")

3. 路径处理

3.1 路径操作

	path := XFile.PathJoin("path", "to", "file.txt")
	norm := XFile.NormalizePath("path/./to/../file.txt")

更多信息请参考模块文档。
*/
package XFile
