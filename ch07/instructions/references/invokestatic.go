package references

import (
	"jvmgo/ch07/instructions/base"
	"jvmgo/ch07/rtda"
	"jvmgo/ch07/rtda/heap"
)

/* invokestatic指令调用静态方法  */
type INVOKE_STATIC struct {
	base.Index16Instruction
}

func (self *INVOKE_STATIC) Execute(frame *rtda.Frame) {
	// 获取常量池
	cp := frame.Method().Class().ConstantPool()
	// 获取方法引用
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	// 解析方法引用, 获取方法
	resolvedMethod := methodRef.ResolvedMethod()
	// 方法必须是静态方法
	if !resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 如果声明方法的类未初始化，需要先初始化类
	//class := resolvedMethod.Class()
	//if !class.InitStarted() {
	//	frame.RevertNextPC()
	//	base.InitClass(frame.Thread(), class)
	//	return
	//}

	base.InvokeMethod(frame, resolvedMethod)
}
