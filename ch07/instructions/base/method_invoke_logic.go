package base

import (
	"jvmgo/ch07/rtda"
	"jvmgo/ch07/rtda/heap"
)

/**
在定位到需要调用的方法之后，Java虚拟机要给这个方法创建一个新的帧并把它推入Java虚拟机栈顶，然后传递参数
*/

func InvokeMethod(invokerFrame *rtda.Frame, method *heap.Method) {
	// 获取栈帧所在的进程
	thread := invokerFrame.Thread()
	// 创建一个新的栈帧
	newFrame := thread.NewFrame(method)
	// 将新创建的栈帧放到栈顶
	thread.PushFrame(newFrame)
	/**
	确认下方法method的参数有几个,.在局部变量表中占用多少位置
	这个数量并不一定等于从Java代码中看到的参数个数，原因有两个：
		第一，long和double类型的参数要占用两个位置。
		第二，对于实例方法，Java编译器会在参数列表的前面添加一个参数，这个隐藏的参数就是this引用
	*/
	argSlotSlot := int(method.ArgSlotCount())
	if argSlotSlot > 0 {
		/**
		有 n 个参数, 则从调用者（invokerFrame）的操作数栈中依此弹出n个变量，将其放在被调用者（newFrame）的局部变量表中，参数传递即可完成
		因为栈是倒放的，所以数据应从后往前放
		*/
		for i := argSlotSlot - 1; i >= 0; i-- {
			// 从调用者（invokerFrame）的操作数栈中弹出变量
			slot := invokerFrame.OperandStack().PopSlot()
			// 将其放在被调用者（newFrame）的局部变量表中
			newFrame.LocalVars().SetSlot(uint((i)), slot)
		}
	}

}
