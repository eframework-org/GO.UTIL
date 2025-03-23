// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XString

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// ToInt 将字符串转换为整数。如果转换失败，则返回 0。
// str 是要转换的字符串。
// 返回转换后的整数。
func ToInt(str string) int {
	itr, _ := strconv.Atoi(str)
	return itr
}

// ToString 将整数转换为字符串。
// itr 是要转换的整数。
// 返回转换后的字符串。
func ToString(itr int) string {
	return strconv.Itoa(itr)
}

// ToFixed 将浮点数转换为具有指定小数位数的字符串。
// float 是要转换的浮点数。
// fixed 是可选的小数位数，默认为2。
// 返回格式化后的字符串。
func ToFixed(float any, fixed ...int) string {
	tfixed := 2
	if len(fixed) > 0 {
		tfixed = fixed[0]
	}
	if tfixed < 0 {
		tfixed = 0
	}
	if float != nil {
		switch float.(type) {
		case float32, float64:
			return fmt.Sprintf("%."+ToString(tfixed)+"f", float)
		}
	}
	return ""
}

// Split 使用分隔符分割字符串。
// str 是要分割的字符串。
// sep 是分隔符。
// 返回分割后的字符串切片。
func Split(str string, sep string) []string {
	return strings.Split(str, sep)
}

// IndexOf 查找子字符串在字符串中第一次出现的位置。
// str 是要搜索的字符串。
// of 是要查找的子字符串。
// 返回子字符串第一次出现的索引，如果未找到则返回 -1。
func IndexOf(str string, of string) int {
	return strings.Index(str, of)
}

// LastIndexOf 查找子字符串在字符串中最后一次出现的位置。
// str 是要搜索的字符串。
// of 是要查找的子字符串。
// 返回子字符串最后一次出现的索引，如果未找到则返回 -1。
func LastIndexOf(str string, of string) int {
	return strings.LastIndex(str, of)
}

// StartsWith 检查字符串是否以指定的子字符串开头。
// str 是要检查的字符串。
// of 是要检查的子字符串。
// 返回字符串是否以指定子字符串开头。
func StartsWith(str string, of string) bool {
	return strings.HasPrefix(str, of)
}

// EndsWith 检查字符串是否以指定的子字符串结尾。
// str 是要检查的字符串。
// of 是要检查的子字符串。
// 返回字符串是否以指定子字符串结尾。
func EndsWith(str string, of string) bool {
	return strings.HasSuffix(str, of)
}

// Contains 检查字符串是否包含指定的子字符串。
// str 是要检查的字符串。
// of 是要查找的子字符串。
// 返回字符串是否包含指定子字符串。
func Contains(str string, of string) bool {
	return strings.Contains(str, of)
}

// IsEmpty 检查字符串是否为空。
// str 是要检查的字符串。
// 返回字符串是否为空。
func IsEmpty(str string) bool {
	return str == ""
}

// Sub 提取字符串中指定范围的子字符串。
// str 是源字符串。
// from 是起始索引。
// to 是结束索引。
// 返回提取的子字符串。
func Sub(str string, from int, to int) string {
	rs := []rune(str)
	length := len(rs)
	if from < 0 || to < 0 || from > to {
		return ""
	}
	if to > length {
		to = length
	}
	return string(rs[from:to])
}

// Replace 替换字符串中所有出现的子字符串。
// str 是要处理的字符串。
// from 是要替换的子字符串。
// to 是替换后的子字符串。
// 返回替换后的字符串。
func Replace(str string, from string, to string) string {
	return strings.ReplaceAll(str, from, to)
}

// Trim 移除字符串两端的空格。
// str 是要处理的字符串。
// 返回处理后的字符串。
func Trim(str string) string {
	return strings.Trim(str, " ")
}

// ToBuffer 将字符串转换为字节切片，不进行内存复制。
// s 是要转换的字符串。
// 返回转换后的字节切片。
func ToBuffer(s string) []byte {
	if s == "" {
		return []byte{}
	}
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// FromBuffer 将字节切片转换为字符串，不进行内存复制。
// b 是要转换的字节切片。
// 返回转换后的字符串。
func FromBuffer(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Format 使用格式说明符格式化字符串。
// format 是格式说明符。
// args 是格式化参数。
// 返回格式化后的字符串。
func Format(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// IEval 定义字符串求值的接口。
type IEval interface {
	// Eval 对输入字符串进行求值。
	// input 是要求值的字符串。
	// 返回求值后的字符串。
	Eval(input string) string
}

// Eval 使用提供的求值器或键值映射对字符串进行求值。
// input 是要求值的字符串。
// sources 是可变数量的求值器(IEval)或map[string]string。
// 返回求值后的字符串。
func Eval(input string, sources ...any) string {
	if input == "" {
		return ""
	}

	result := input
	for _, source := range sources {
		if source == nil {
			continue
		}

		switch src := source.(type) {
		case IEval:
			result = src.Eval(result)
		case map[string]string:
			for key, value := range src {
				result = strings.ReplaceAll(result, key, value)
			}
		}
	}

	return result
}

// Random 生成随机字符串。
// format 参数指定格式("D"表示36位带连字符，"N"表示32位无连字符)，
// 返回生成的随机字符串。
//
// 示例：
//
//	random := XString.Random()      // 返回："c9a0cad5e9624b3b8e07f5df0e5c1bbc"
//	random := XString.Random("D")   // 返回："c9a0cad5-e962-4b3b-8e07-f5df0e5c1bbc"
func Random(format ...string) string {
	// Generate 16 random bytes
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	// Set version (4) and variant (2) bits
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 2

	var guid string
	// Format the UUID
	if len(format) > 0 {
		if format[0] == "D" {
			// Format with hyphens (8-4-4-4-12)
			guid = fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
				b[0:4],
				b[4:6],
				b[6:8],
				b[8:10],
				b[10:16])
		}
	}
	if guid == "" {
		// Format without hyphens (32 chars)
		guid = fmt.Sprintf("%08x%04x%04x%04x%012x",
			b[0:4],
			b[4:6],
			b[6:8],
			b[8:10],
			b[10:16])
	}

	return guid
}

// DES加密向量
var rgbIV = []byte{0x7B, 0x4A, 0xF3, 0x91, 0xE5, 0xD2, 0x8C, 0x6F}

// Encrypt 使用DES算法加密数据。
// data 是要加密的数据。
// key 是可选的加密密钥，默认使用rgbIV。
// 返回base64编码的加密字符串。
func Encrypt[T string | []byte](data T, key ...string) string {
	// 检查并转换输入数据
	var content []byte
	switch v := any(data).(type) {
	case string:
		if v == "" {
			return ""
		}
		content = []byte(v)
	case []byte:
		if len(v) == 0 {
			return ""
		}
		content = v
	default:
		return ""
	}

	// 准备加密密钥
	rgb := rgbIV
	if len(key) > 0 && key[0] != "" {
		rgb = padKey([]byte(key[0]))
	}

	// 创建DES加密块
	block, err := des.NewCipher(rgb)
	if err != nil {
		return ""
	}

	// 填充数据并加密
	content = pkcs7Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	cipher.NewCBCEncrypter(block, rgbIV).CryptBlocks(crypted, content)

	return base64.StdEncoding.EncodeToString(crypted)
}

// Decrypt 解密DES加密的数据。
// data 是base64编码的加密数据。
// key 是可选的解密密钥，默认使用rgbIV。
// 返回解密后的数据。
func Decrypt[T string | []byte](data string, key ...string) T {
	var zero T
	if data == "" {
		return zero
	}

	// 解码base64数据
	crypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return zero
	}

	// 准备解密密钥
	rgb := rgbIV
	if len(key) > 0 && key[0] != "" {
		rgb = padKey([]byte(key[0]))
	}

	// 创建DES解密块
	block, err := des.NewCipher(rgb)
	if err != nil {
		return zero
	}

	// 解密数据
	origData := make([]byte, len(crypted))
	cipher.NewCBCDecrypter(block, rgbIV).CryptBlocks(origData, crypted)
	origData = pkcs7UnPadding(origData)

	// 根据返回类型转换结果
	switch any(zero).(type) {
	case string:
		return any(string(origData)).(T)
	case []byte:
		return any(origData).(T)
	default:
		return zero
	}
}

// padKey 确保密钥长度为8字节。
// key 是要处理的密钥。
// 返回处理后的8字节密钥。
func padKey(key []byte) []byte {
	if len(key) > 8 {
		return key[:8]
	}
	return append(key, bytes.Repeat([]byte{0}, 8-len(key))...)
}

// pkcs7Padding 对数据进行PKCS7填充。
// data 是要填充的数据。
// blockSize 是块大小。
// 返回填充后的数据。
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7UnPadding 移除数据的PKCS7填充。
// data 是要移除填充的数据。
// 返回移除填充后的数据。
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}
	return data[:(length - unpadding)]
}
