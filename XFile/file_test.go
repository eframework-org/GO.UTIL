// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XFile

import (
	"os"
	"testing"
)

// 检查文件是否存在
func TestHasFile(t *testing.T) {
	testFile := "testfile.txt"
	os.WriteFile(testFile, []byte("test content"), 0644)
	defer os.Remove(testFile)

	if !HasFile(testFile) {
		t.Errorf("Expected file to exist: %s", testFile)
	}

	if HasFile("nonexistentfile.txt") {
		t.Errorf("Expected file to not exist: nonexistentfile.txt")
	}
}

// 测试打开文件并读取其内容
func TestOpenFile(t *testing.T) {
	testFile := "testfile.txt"
	content := []byte("test content")
	os.WriteFile(testFile, content, 0644)
	defer os.Remove(testFile)

	data := OpenFile(testFile)
	if string(data) != string(content) {
		t.Errorf("Expected file content: %s, got: %s", content, data)
	}

	data = OpenFile("nonexistentfile.txt")
	if data != nil {
		t.Errorf("Expected nil for nonexistent file, got: %s", data)
	}
}

// 测试将内容保存到文件
func TestSaveFile(t *testing.T) {
	testFile := "testfile.txt"
	content := []byte("test content")
	err := SaveFile(testFile, content, 0644)
	if err != nil {
		t.Errorf("Failed to save file: %v", err)
	}
	defer os.Remove(testFile)

	data, _ := os.ReadFile(testFile)
	if string(data) != string(content) {
		t.Errorf("Expected file content: %s, got: %s", content, data)
	}
}

// 测试文件的删除
func TestDeleteFile(t *testing.T) {
	testFile := "testfile.txt"
	os.WriteFile(testFile, []byte("test content"), 0644)

	if !DeleteFile(testFile) {
		t.Errorf("Failed to delete file: %s", testFile)
	}

	if HasFile(testFile) {
		t.Errorf("Expected file to be deleted: %s", testFile)
	}
}

// 测试打开文本文件并读取其内容
func TestOpenText(t *testing.T) {
	testFile := "testfile.txt"
	content := "test content"
	os.WriteFile(testFile, []byte(content), 0644)
	defer os.Remove(testFile)

	data := OpenText(testFile)
	if data != content {
		t.Errorf("Expected file content: %s, got: %s", content, data)
	}

	data = OpenText("nonexistentfile.txt")
	if data != "" {
		t.Errorf("Expected empty string for nonexistent file, got: %s", data)
	}
}

// 测试将文本内容保存到文件
func TestSaveText(t *testing.T) {
	testFile := "testfile.txt"
	content := "test content"
	err := SaveText(testFile, content, 0644)
	if err != nil {
		t.Errorf("Failed to save file: %v", err)
	}
	defer os.Remove(testFile)

	data, _ := os.ReadFile(testFile)
	if string(data) != content {
		t.Errorf("Expected file content: %s, got: %s", content, data)
	}
}

// 检查目录是否存在
func TestHasDirectory(t *testing.T) {
	testDir := "testdir"
	os.Mkdir(testDir, os.ModePerm)
	defer os.RemoveAll(testDir)

	if !HasDirectory(testDir) {
		t.Errorf("Expected directory to exist: %s", testDir)
	}

	if HasDirectory("nonexistentdir") {
		t.Errorf("Expected directory to not exist: nonexistentdir")
	}
}

// 测试创建新目录并处理已存在的目录
func TestCreateDirectory(t *testing.T) {
	testDir := "testdir"

	// Clean up before test
	if HasDirectory(testDir) {
		os.RemoveAll(testDir)
	}

	// Test creating a new directory
	if !CreateDirectory(testDir) {
		t.Errorf("Failed to create directory: %s", testDir)
	}

	// Check if the directory was created
	if !HasDirectory(testDir) {
		t.Errorf("Directory was not created: %s", testDir)
	}

	// Test creating an existing directory
	if !CreateDirectory(testDir) {
		t.Errorf("Failed to create existing directory: %s", testDir)
	}

	// Clean up after test
	os.RemoveAll(testDir)
}

// 测试从路径中提取目录名称
func TestDirectoryName(t *testing.T) {
	path := "path/to/file.txt"
	expected := "path/to"
	if dir := DirectoryName(path); dir != expected {
		t.Errorf("Expected directory: %s, got: %s", expected, dir)
	}
}

// 测试将多个路径段连接成一个单一路径
func TestPathJoin(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{""}, ""},
		{[]string{"a", "b", "c"}, "a/b/c"},
		{[]string{"a/", "/b", "c"}, "a/b/c"},
		{[]string{"a", "", "c"}, "a/c"},
		{[]string{"a", ".", "c"}, "a/c"},
		{[]string{"a", "..", "c"}, "c"},
		{[]string{"a", "b", "../c"}, "a/c"},
		{[]string{"a", "b", "c/.."}, "a/b"},
		{[]string{"a", "b", "c/../.."}, "a"},
		{[]string{"a", "b", "c/../../.."}, "."},
	}

	for _, test := range tests {
		result := PathJoin(test.input...)
		if result != test.expected {
			t.Errorf("PathJoin(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

// 测试规范化各种路径格式
func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"./", "."},
		{"../", ""},
		{"a/b/c", "a/b/c"},
		{"a\\b\\c", "a/b/c"},
		{"a/./b/./c", "a/b/c"},
		{"a/b/../c", "a/c"},
		{"a/b/c/..", "a/b"},
		{"a/b/c/../..", "a"},
		{"a/b/c/../../..", ""},
		{"a/b/c/../../../..", ""},
	}

	for _, test := range tests {
		result := NormalizePath(test.input)
		if result != test.expected {
			t.Errorf("NormalizePath(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
