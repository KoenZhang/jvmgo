package references

import (
	"jvmgo/ch06/instructions/base"
	"jvmgo/ch06/rtda"
	"jvmgo/ch06/rtda/heap"
)

/**
ldc系列指令 从运行时常量池中加载常量值，并把它推入操作数栈。ldc系列指令属于常量类指令，共3条:
1. ldc 和 ldc_w 指令用于加载 int、float 和 字符串常量，java.lang.Class实例 或者 MethodType 和 MethodHandle实例, ldc和ldc_w指令的区别仅在于操作数的宽度
2. ldc2_w 指令用于加载 long 和 double 常量。
*/

type LDC struct {
	base.Index8Instruction
}

type LDC_W struct {
	base.Index16Instruction
}

type LDC2_W struct {
	base.Index16Instruction
}

func (self *LDC) Execute(frame *rtda.Frame) {
	_ldc(frame, self.Index)
}

func (self *LDC_W) Execute(frame *rtda.Frame) {
	_ldc(frame, self.Index)
}

func (self *LDC2_W) Execute(frame *rtda.Frame) {
	// 获取操作数栈
	stack := frame.OperandStack()
	cp := frame.Method().Class().ConstantPool()
	c := cp.GetConstant(self.Index)
	switch c.(type) {
	case int64:
		stack.PushLong(c.(int64))
	case float64:
		stack.PushDouble(c.(float64))
	default:
		panic("java.lang.ClassFormatError")
	}
}

// 从运行时常量池中加载常量值，并把它推入操作数栈
func _ldc(frame *rtda.Frame, index uint) {
	// 操作数栈
	stack := frame.OperandStack()
	// 运行时常量池
	cp := frame.Method().Class().ConstantPool()
	// 获取常量值
	c := cp.GetConstant(index)

	// 推入操作数栈
	switch c.(type) {
	case int32:
		stack.PushInt(c.(int32))
	case float32:
		stack.PushFloat(c.(float32))
	case string:
		// todo
	case *heap.ClassRef:
		// todo
	default:
		panic("todo: ldc!")
	}
}
