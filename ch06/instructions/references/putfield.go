package references

import (
	"jvmgo/ch06/instructions/base"
	"jvmgo/ch06/rtda"
	"jvmgo/ch06/rtda/heap"
)

/**
putfield指令给 实例变量赋值 , 它需要三个操作数：
1. 第一个操作数是常量池索引，16位
2. 第一个操作数变量值，从操作数栈中弹出
3. 第三个操作数是对象引用，从操作数栈中弹出
*/

type PUT_FIELD struct {
	// 常量池索引
	base.Index16Instruction
}

func (self *PUT_FIELD) Execute(frame *rtda.Frame) {
	// 获取当前方法
	currentMethod := frame.Method()
	// 根据当前方法获取当前类
	currentClass := currentMethod.Class()
	// 根据当前类获取常量池
	cp := currentClass.ConstantPool()
	// 根据常量池获取字段引用
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
	// 解析字段引用获取字段
	field := fieldRef.ResolvedField()

	// 解析后的字段必须是实例字段，否则抛出IncompatibleClassChangeError
	if field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 如果是final字段，则只能在构造函数中<init>初始化，否则抛出IllegalAccessError
	if field.IsFinal() {
		if currentClass != field.Class() || currentMethod.Name() != "<init>" {
			panic("java.lang.IllegalAccessError")
		}
	}

	// 先根据字段类型从操作数栈中弹出相应的变量值，然后弹出对象引用。如果引用是null，需要抛出著名的空指针异常（NullPointerException），否则通过引用给实例变量赋值
	// 获取字段描述符
	descriptor := field.Descriptor()
	// 获取实例字段所在位置
	slotId := field.SlotId()
	// 操作数栈
	stack := frame.OperandStack()
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		// 根据字段类型从操作数栈中弹出相应的变量值
		val := stack.PopInt()
		// 弹出对象引用
		ref := stack.PopRef()
		// 如果引用是null，需要抛出著名的空指针异常（NullPointerException）
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetInt(slotId, val)
	case 'F':
		val := stack.PopFloat()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetFloat(slotId, val)
	case 'J':
		val := stack.PopLong()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetLong(slotId, val)
	case 'D':
		val := stack.PopDouble()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetDouble(slotId, val)
	case 'L', '[':
		val := stack.PopRef()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetRef(slotId, val)
	default:
		// todo
	}
}
