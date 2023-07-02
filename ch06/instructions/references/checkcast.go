package references

import (
	"jvmgo/ch06/instructions/base"
	"jvmgo/ch06/rtda"
	"jvmgo/ch06/rtda/heap"
)

/**
checkcast指令和instanceof指令很像，区别在于：
instanceof指令会改变操作数栈（弹出对象引用，推入判断结果）
checkcast则不改变操作数栈（如果判断失败，直接抛出ClassCastException异常）
*/

type CHECK_CAST struct {
	base.Index16Instruction
}

func (self *CHECK_CAST) Execute(frame *rtda.Frame) {
	// 操作数栈
	stack := frame.OperandStack()
	// 弹出引用
	ref := stack.PopRef()
	// 再将引用放回操作数栈
	stack.PushRef(ref)
	// 如果引用是 null，则指令执行结束。也就是说，null引用可以转换成任何类型
	if ref == nil {
		return
	}

	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	if !ref.IsInstanceOf(class) {
		panic("java.lang.ClassCastException")
	}
}
