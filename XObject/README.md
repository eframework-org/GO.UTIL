# XObject

[![Reference](https://pkg.go.dev/badge/github.com/eframework-org/GO.UTIL/XObject.svg)](https://pkg.go.dev/github.com/eframework-org/GO.UTIL/XObject)
[![Release](https://img.shields.io/github/v/tag/eframework-org/GO.UTIL)](https://github.com/eframework-org/GO.UTIL/tags)
[![Report](https://goreportcard.com/badge/github.com/eframework-org/GO.UTIL)](https://goreportcard.com/report/github.com/eframework-org/GO.UTIL)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-Explore-blue)](https://deepwiki.com/eframework-org/GO.UTIL)

XObject 提供了 Go 语言面向对象编程的增强支持，包括对象构造、实例管理和序列化功能。

## 功能特性

- 对象构造器：支持泛型的对象构造和初始化
- 实例指针：支持对象实例的自引用
- 参数化构造：支持最多三个参数的构造函数
- JSON 序列化：支持对象与 JSON 的相互转换

## 使用手册

### 1. 对象构造

#### 1.1 基础构造
```go
// 定义结构体
type MyStruct struct {
    XObject.ICtor // 实现无参构造
    Name string
}

// 实现构造函数
func (ms *MyStruct) Ctor(obj any) {
    ms.Name = "默认名称"
}

// 创建实例
obj := XObject.New[MyStruct]()
```

#### 1.2 参数化构造
```go
// 定义结构体
type MyStruct struct {
    XObject.ICtorT1[string] // 实现带参构造
    Name string
}

// 实现构造函数
func (ms *MyStruct) CtorT1(obj any, name string) {
    ms.Name = name
}

// 创建实例
obj := XObject.NewT1[MyStruct, string]("张三")
```

#### 1.3 实例自引用
```go
// 定义结构体
type MyStruct struct {
    XObject.IThis[MyStruct] // 实现自引用
    this *MyStruct
}

// 实现自引用方法
func (ms *MyStruct) This() *MyStruct {
    return ms.this
}

// 在构造函数中初始化
func (ms *MyStruct) Ctor(obj any) {
    ms.this = obj.(*MyStruct)
}
```

### 2. 序列化

#### 2.1 JSON 转换
```go
// 对象转 JSON
json, err := XObject.ToJson(obj)
json, err := XObject.ToJson(obj, true) // 格式化输出

// JSON 转对象
err := XObject.FromJson(json, &obj)
```

#### 2.2 字节转换
```go
// 对象转字节数组
bytes, err := XObject.ToByte(obj)

// 字节数组转对象
err := XObject.FromByte(bytes, &obj)
```

## 常见问题

### 1. 为什么需要实现构造器接口？
构造器接口提供了对象初始化的标准方式：
- 确保对象在创建时正确初始化
- 支持依赖注入和参数化构造
- 便于实现工厂模式和对象池
- 统一的对象生命周期管理

### 2. 如何选择合适的构造器？
根据初始化参数的数量选择：
- ICtor：无参数初始化
  - 适用于简单对象
  - 默认值初始化
  - 单例模式
- ICtorT1：一个参数
  - 基本配置初始化
  - 依赖注入
- ICtorT2：两个参数
  - 复杂配置初始化
  - 多依赖注入
- ICtorT3：三个参数
  - 完整状态初始化
  - 完整依赖注入

### 3. JSON 序列化失败的常见原因
1. 字段可见性：
   - 确保需要序列化的字段首字母大写
   - 使用 `json` 标签控制序列化行为

2. 类型兼容：
   - 确保字段类型支持序列化
   - 注意数字类型的精度
   - 处理时间类型的格式

3. 循环引用：
   - 避免对象间的循环引用
   - 使用指针打破循环依赖
   - 考虑使用自定义 MarshalJSON

### 4. 性能优化建议

1. 对象创建：
   - 使用对象池复用实例
   - 避免频繁创建临时对象
   - 合理使用指针和值类型

2. 序列化优化：
   - 使用字节数组而不是字符串
   - 避免不必要的格式化
   - 考虑使用二进制序列化

3. 内存管理：
   - 及时释放不用的对象
   - 避免大对象的深拷贝
   - 使用 sync.Pool 管理临时对象

更多问题，请查阅[问题反馈](../CONTRIBUTING.md#问题反馈)。

## 项目信息

- [更新记录](../CHANGELOG.md)
- [贡献指南](../CONTRIBUTING.md)
- [许可证](../LICENSE)
