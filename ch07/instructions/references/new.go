package references

import (
	"jvmgo/ch07/instructions/base"
	"jvmgo/ch07/rtda"
	"jvmgo/ch07/rtda/heap"
)

// new指令专门用来创建类实例: new Object

type NEW struct {
	// new指令的操作数是一个uint16索引，来自字节码
	base.Index16Instruction
}

/**
new指令的操作数是一个uint16索引，来自字节码。
通过这个索引，可以从当前类的运行时常量池中找到一个类符号引用。
解析这个类符号引用，拿到类数据，然后创建对象，并把对象引用推入栈顶，new指令的工作就完成了
*/
func (self *NEW) Execute(frame *rtda.Frame) {
	// 获取类对应的常量池
	cp := frame.Method().Class().ConstantPool()
	// 从当前类的运行时常量池中找到一个类符号引用
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	// 接口和抽象类都不能实例化
	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}
	// 创建对象
	ref := class.NewObject()
	// 将对象引用推入栈顶
	frame.OperandStack().PushRef(ref)
}
