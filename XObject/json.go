// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XObject

import (
	"bytes"
	"encoding/json"

	"github.com/eframework-org/GO.UTIL/XString"
)

// encode 将对象编码为字节数组。
// 如果对象是字符串类型，直接返回其字节表示。
// 支持通过 pretty 参数控制是否格式化输出。
func encode(v any, pretty ...bool) ([]byte, error) {
	if str, ok := v.(string); ok {
		return XString.ToBuffer(str), nil
	}

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	encoder.SetEscapeHTML(false)
	if len(pretty) > 0 && pretty[0] {
		encoder.SetIndent("", "    ")
	}
	if e := encoder.Encode(v); e != nil {
		return nil, e
	}
	bytes := buff.Bytes()
	size := len(bytes)
	if size > 0 && bytes[size-1] == '\n' {
		bytes = bytes[:size-1]
	}

	return bytes, nil
}

// ToJson 将对象转换为 JSON 字符串。
// 支持通过 pretty 参数控制是否格式化输出。
// 如果转换失败，返回空字符串和错误信息。
func ToJson(v any, pretty ...bool) (string, error) {
	b, e := encode(v, pretty...)
	if e != nil {
		return "", e
	}
	return XString.FromBuffer(b), nil
}

// FromJson 将 JSON 字符串转换为指定类型的对象。
// 如果转换失败，返回错误信息。
func FromJson[T any](data string, obj T) error {
	if e := json.Unmarshal(XString.ToBuffer(data), obj); e != nil {
		return e
	}
	return nil
}

// ToByte 将对象转换为字节数组。
// 如果对象是字符串类型，直接返回其字节表示。
// 如果转换失败，返回 nil 和错误信息。
func ToByte(v any) ([]byte, error) {
	b, e := encode(v)
	if e != nil {
		return nil, e
	}
	return b, nil
}

// FromByte 将字节数组转换为指定类型的对象。
// 如果转换失败，返回错误信息。
func FromByte[T any](data []byte, obj T) error {
	if e := json.Unmarshal(data, obj); e != nil {
		return e
	}
	return nil
}
