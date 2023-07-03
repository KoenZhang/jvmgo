package references

import "jvmgo/ch07/instructions/base"
import "jvmgo/ch07/rtda"

/**
因为对象是需要初始化的，所以每个类都至少有一个构造函数。
即使用户自己不定义，编译器也会自动生成一个默认构造函数。
在创建类实例时，编译器会在 new指令 的后面加入 invokespecial 指令来调用构造函数初始化对象
*/
type INVOKE_SPECIAL struct{ base.Index16Instruction }

// hack!
func (self *INVOKE_SPECIAL) Execute(frame *rtda.Frame) {
	frame.OperandStack().PopRef()
}
