package control

import (
	"jvmgo/ch07/instructions/base"
	"jvmgo/ch07/rtda"
)

/**
方法执行完毕之后，需要把结果返回给调用方，这一工作由返回指令完成。返回指令属于控制类指令，一共有6条:
	1. return	指令用于没有返回值的情况
	2. areturn	指令用于返回 引用 类型的值
	3. ireturn 	指令用于返回 int 类型的值
	4. lreturn	指令用于返回 long 类型的值
	5. freturn	指令用于返回 float 类型的值
	6. dreturn	指令用于返回 double 类型的值
*/

// 用于没有返回值的情况
type RETURN struct {
	base.NoOperandsInstruction
}

// return指令比较简单，只要把当前帧从Java虚拟机栈中弹出即可
func (self *RETURN) Execute(frame *rtda.Frame) {
	frame.Thread().PopFrame()
}

// 用于返回 引用 类型的值
type ARETURN struct {
	base.NoOperandsInstruction
}

func (self *ARETURN) Execute(frame *rtda.Frame) {
	// 根据栈帧获取当前进程
	thread := frame.Thread()
	// 弹出栈帧，根据当前进程获取被调用者，现在被调用者已经执行完毕
	currentFrame := thread.PopFrame()
	// 获取被调用者的返回值
	retVal := currentFrame.OperandStack().PopRef()

	// 现在进程最新的方法是调用者，这时直接获取即可，不用弹出栈帧
	invokerFrame := thread.TopFrame()
	// 将被调用者执行后的返回值压入调用者操作数栈栈顶
	invokerFrame.OperandStack().PushRef(retVal)
}

// 用于返回 double 类型的值
type DRETURN struct {
	base.NoOperandsInstruction
}

func (self *DRETURN) Execute(frame *rtda.Frame) {
	// 根据栈帧获取当前进程
	thread := frame.Thread()
	// 弹出栈帧，根据当前进程获取被调用者，现在被调用者已经执行完毕
	currentFrame := thread.PopFrame()
	// 获取被调用者的返回值
	retVal := currentFrame.OperandStack().PopDouble()

	// 现在进程最新的方法是调用者，这时直接获取即可，不用弹出栈帧
	invokerFrame := thread.TopFrame()
	// 将被调用者执行后的返回值压入调用者操作数栈栈顶
	invokerFrame.OperandStack().PushDouble(retVal)
}

// 用于返回 float 类型的值
type FRETURN struct {
	base.NoOperandsInstruction
}

func (self *FRETURN) Execute(frame *rtda.Frame) {
	// 根据栈帧获取当前进程
	thread := frame.Thread()
	// 弹出栈帧，根据当前进程获取被调用者，现在被调用者已经执行完毕
	currentFrame := thread.PopFrame()
	// 获取被调用者的返回值
	retVal := currentFrame.OperandStack().PopFloat()

	// 现在进程最新的方法是调用者，这时直接获取即可，不用弹出栈帧
	invokerFrame := thread.TopFrame()
	// 将被调用者执行后的返回值压入调用者操作数栈栈顶
	invokerFrame.OperandStack().PushFloat(retVal)
}

// 用于返回 int 类型的值
type IRETURN struct {
	base.NoOperandsInstruction
}

func (self *IRETURN) Execute(frame *rtda.Frame) {
	// 根据栈帧获取当前进程
	thread := frame.Thread()
	// 弹出栈帧，根据当前进程获取被调用者，现在被调用者已经执行完毕
	currentFrame := thread.PopFrame()
	// 获取被调用者的返回值
	retVal := currentFrame.OperandStack().PopInt()

	// 现在进程最新的方法是调用者，这时直接获取即可，不用弹出栈帧
	invokerFrame := thread.TopFrame()
	// 将被调用者执行后的返回值压入调用者操作数栈栈顶
	invokerFrame.OperandStack().PushInt(retVal)
}

// 用于返回 long 类型的值
type LRETURN struct {
	base.NoOperandsInstruction
}

func (self *LRETURN) Execute(frame *rtda.Frame) {
	// 根据栈帧获取当前进程
	thread := frame.Thread()
	// 弹出栈帧，根据当前进程获取被调用者，现在被调用者已经执行完毕
	currentFrame := thread.PopFrame()
	// 获取被调用者的返回值
	retVal := currentFrame.OperandStack().PopLong()

	// 现在进程最新的方法是调用者，这时直接获取即可，不用弹出栈帧
	invokerFrame := thread.TopFrame()
	// 将被调用者执行后的返回值压入调用者操作数栈栈顶
	invokerFrame.OperandStack().PushLong(retVal)
}
