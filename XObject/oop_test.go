package XObject

import (
	"testing"
)

type TestStruct struct {
	ICtor

	Inited bool
}

func (ts *TestStruct) Ctor(obj any) { ts.Inited = true }

func TestNew(t *testing.T) {
	obj := New[TestStruct]()
	if obj == nil {
		t.Error("Expected non-nil object")
	}
}

type TestStructT1 struct {
	ICtorT1[int]

	Value1 int
}

func (tst1 *TestStructT1) CtorT1(obj any, arg1 int) {
	tst1.Value1 = arg1
}

func TestNewT1(t *testing.T) {
	obj := NewT1[TestStructT1](42)
	if obj == nil {
		t.Error("Expected non-nil object")
	}
	if obj.Value1 != 42 {
		t.Errorf("Expected Value1 to be 42, got %d", obj.Value1)
	}
}

type TestStructT2 struct {
	ICtorT2[int, string]

	Value1 int
	Value2 string
}

func (tst2 *TestStructT2) CtorT2(obj any, arg1 int, arg2 string) {
	tst2.Value1 = arg1
	tst2.Value2 = arg2
}

func TestNewT2(t *testing.T) {
	obj := NewT2[TestStructT2](42, "hello")
	if obj == nil {
		t.Error("Expected non-nil object")
	}
	if obj.Value1 != 42 {
		t.Errorf("Expected Value1 to be 42, got %d", obj.Value1)
	}
	if obj.Value2 != "hello" {
		t.Errorf("Expected Value2 to be 'hello', got '%s'", obj.Value2)
	}
}

type TestStructT3 struct {
	ICtorT3[int, string, bool]

	Value1 int
	Value2 string
	Value3 bool
}

func (tst3 *TestStructT3) CtorT3(obj any, arg1 int, arg2 string, arg3 bool) {
	tst3.Value1 = arg1
	tst3.Value2 = arg2
	tst3.Value3 = arg3
}

func TestNewT3(t *testing.T) {
	obj := NewT3[TestStructT3](42, "hello", true)
	if obj == nil {
		t.Error("Expected non-nil object")
	}
	if obj.Value1 != 42 {
		t.Errorf("Expected Value1 to be 42, got %d", obj.Value1)
	}
	if obj.Value2 != "hello" {
		t.Errorf("Expected Value2 to be 'hello', got '%s'", obj.Value2)
	}
	if obj.Value3 != true {
		t.Errorf("Expected Value3 to be true, got %v", obj.Value3)
	}
}

type TestBaseStruct struct {
	ICtor

	Value int
}

func (tbs *TestBaseStruct) Ctor(obj any) { tbs.Value = 1 }

type TestSubStruct struct {
	ICtor
	TestBaseStruct

	Value int
}

func (tss *TestSubStruct) Ctor(obj any) {
	tss.TestBaseStruct.Ctor(obj)
	tss.Value = 1
}

func TestInherit(t *testing.T) {
	obj := New[TestSubStruct]()
	if obj == nil {
		t.Error("Expected non-nil object")
	}
	if obj.TestBaseStruct.Value != 1 {
		t.Errorf("Expected TestBaseStruct.Value to be 1, got %d", obj.TestBaseStruct.Value)
	}
	if obj.Value != 1 {
		t.Errorf("Expected TestSubStruct.Value to be 1, got %d", obj.Value)
	}
}

type TestThisStruct struct {
	ICtor
	IThis[TestBaseStruct]

	this *TestThisStruct
}

func (tt *TestThisStruct) This() *TestThisStruct { return tt.this }

func (tt *TestThisStruct) Ctor(obj any) {
	tt.this = obj.(*TestThisStruct)
}

func TestThis(t *testing.T) {
	obj := New[TestThisStruct]()
	if obj == nil {
		t.Error("Expected non-nil object")
	}
	if obj != obj.This() {
		t.Errorf("Expected TestThisStruct.This() to be obj, got %v", obj.This())
	}
}
