package references

import (
	"jvmgo/ch07/instructions/base"
	"jvmgo/ch07/rtda"
	"jvmgo/ch07/rtda/heap"
)

/*
调用需要动态绑定的方法，针对接口类型的引用调用方法，就使用invokeinterface指令
在字节码中，invokeinterface指令的操作码后面跟着4字节而非2字节：
	1. 前两字节的含义和其他指令相同，是个uint16运行时常量池索引。
	2. 第3字节的值是给方法传递参数需要的slot数，其含义和给Method结构体定义的argSlotCount字段相同。正如我们所知，这个数是可以根据方法描述符计算出来的，它的存在仅仅是因为历史原因。
	3. 第4字节是留给Oracle的某些Java虚拟机实现用的，它的值必须是0。该字节的存在是为了保证Java虚拟机可以向后兼容
*/
type INVOKE_INTERFACE struct {
	index uint
}

func (self *INVOKE_INTERFACE) FetchOperands(reader *base.BytecodeReader) {
	self.index = uint(reader.ReadUint16()) // 运行时常量池索引
	reader.ReadUint8()                     // count		slot数
	reader.ReadUint8()                     // must be 0 Oracle的某些Java虚拟机实现用的，它的值必须是0。该字节的存在是为了保证Java虚拟机可以向后兼容
}

func (self *INVOKE_INTERFACE) Execute(frame *rtda.Frame) {
	// 获取运行时常量池
	cp := frame.Method().Class().ConstantPool()
	// 获取方法符号引用
	methodRef := cp.GetConstant(self.index).(*heap.InterfaceMethodRef)
	// 解析方法符号引用，获取方法
	resolvedMethod := methodRef.ResolvedInterfaceMethod()
	// 解析后的方法为 静态方法 或者 私有方法，则抛出异常
	if resolvedMethod.IsStatic() || resolvedMethod.IsPrivate() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 从操作数栈中获取this引用，如果引用是null，则抛出NullPointerException异常
	ref := frame.OperandStack().GetRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if ref == nil {
		panic("java.lang.NullPointerException")
	}
	// 如果引用所指对象的类没有实现解析出来的接口，则抛出IncompatibleClassChangeError异常
	if !ref.Class().IsImplements(methodRef.ResolvedClass()) {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 查找最终要调用的方法
	methodToBeInvoked := heap.LookupMethodInClass(ref.Class(), methodRef.Name(), methodRef.Descriptor())
	// 如果找不到，或者找到的方法是抽象的，则抛出Abstract-MethodError异常。如果找到的方法不是public，则抛出IllegalAccessError异常
	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}
	if !methodToBeInvoked.IsPublic() {
		panic("java.lang.IllegalAccessError")
	}

	// 正常调用方法
	base.InvokeMethod(frame, methodToBeInvoked)

}
