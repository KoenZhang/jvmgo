package references

import (
	"jvmgo/ch06/instructions/base"
	"jvmgo/ch06/rtda"
	"jvmgo/ch06/rtda/heap"
)

/** instanceof指令判断对象是否是某个类的实例（或者对象的类是否实现了某个接口）
instanceof指令需要两个操作数:
1. 第一个操作数是uint16索引，从方法的字节码中获取，通过这个索引可以从当前类的运行时常量池中找到一个类符号引用
2. 第二个操作数是对象引用，从操作数栈中弹出
*/
type INSTANCE_OF struct {
	base.Index16Instruction
}

/**
instanceof指令需要两个操作数:
1. 第一个操作数是uint16索引，从方法的字节码中获取，通过这个索引可以从当前类的运行时常量池中找到一个类符号引用
2. 第二个操作数是对象引用，从操作数栈中弹出
*/
func (self *INSTANCE_OF) Execute(frame rtda.Frame) {
	stack := frame.OperandStack()
	ref := stack.PopRef() // 对象引用
	// 先弹出对象引用，如果是null，则把 0 推入操作数栈
	if ref == nil {
		stack.PushInt(0)
	}

	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolverdClass()
	if ref.IsInstanceOf(class) {
		stack.PushInt(1)
	} else {
		stack.PushInt(0)
	}
}
