// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
XString 提供了高效的字符串处理工具，支持转换、操作、求值和加密等功能。

功能特性

  - 字符串转换：支持整数、浮点数与字符串之间的转换，提供零拷贝的字节切片转换
  - 字符串操作：提供分割、查找、替换、修剪等基本操作
  - 字符串求值：支持可扩展的字符串求值系统，支持自定义求值器
  - 加密解密：支持 DES 加密的字符串安全处理，支持自定义密钥
  - 随机字符串：支持生成指定格式和长度的随机字符串

使用手册

1. 字符串转换

1.1 整数转换

	// 字符串转整数
	num := XString.ToInt("123")     // 返回：123

	// 整数转字符串
	str := XString.ToString(456)    // 返回："456"

1.2 浮点数格式化

	// 格式化浮点数，指定小数位数
	fixed := XString.ToFixed(3.14159, 2)  // 返回："3.14"

1.3 字节切片转换

	// 字符串转字节切片（零拷贝）
	buf := XString.ToBuffer("hello")

	// 字节切片转字符串（零拷贝）
	str := XString.FromBuffer(buf)

2. 字符串操作

2.1 基本操作

	// 分割字符串
	parts := XString.Split("a,b,c", ",")  // 返回：["a", "b", "c"]

	// 查找子串
	pos := XString.IndexOf("hello", "l")   // 返回：2
	last := XString.LastIndexOf("hello", "l") // 返回：3

	// 提取子串
	sub := XString.Sub("hello", 0, 2)      // 返回："he"

2.2 判断操作

	// 检查前缀后缀
	hasPrefix := XString.StartsWith("hello", "he") // 返回：true
	hasSuffix := XString.EndsWith("hello", "lo")   // 返回：true

	// 包含检查
	contains := XString.Contains("hello", "ll")    // 返回：true

	// 空值检查
	isEmpty := XString.IsEmpty("")               // 返回：true

3. 字符串求值

3.1 Map 方式求值

	data := map[string]string{
		"${name}": "John",
		"${age}": "30",
	}
	result := XString.Eval("Name: ${name}, Age: ${age}", data)
	// 返回："Name: John, Age: 30"

3.2 自定义求值器

	type CustomEval struct{}

	func (e *CustomEval) Eval(input string) string {
		// 自定义求值逻辑
		return input
	}

	evaluator := &CustomEval{}
	result := XString.Eval("some text", evaluator)

4. 加密解密

4.1 字符串加密

	// 使用默认密钥加密
	encrypted := XString.Encrypt("sensitive data")

	// 使用自定义密钥加密
	encrypted := XString.Encrypt("sensitive data", "mykey")

4.2 字符串解密

	// 使用默认密钥解密
	decrypted := XString.Decrypt(encrypted)

	// 使用自定义密钥解密
	decrypted := XString.Decrypt(encrypted, "mykey")

5. 随机字符串

5.1 生成随机字符串

	// 生成32位无连字符的随机字符串
	random := XString.Random("N")  // 返回：如 "c9a0cad5e9624b3b8e07f5df0e5c1bbc"

	// 生成36位带连字符的随机字符串
	random := XString.Random("D")  // 返回：如 "c9a0cad5-e962-4b3b-8e07-f5df0e5c1bbc"

	// 生成指定长度的随机字符串
	random := XString.Random("N", 16)  // 返回：如 "c9a0cad5e9624b3b"

更多信息请参考模块文档。
*/
package XString
