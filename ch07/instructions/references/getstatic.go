package references

import (
	"jvmgo/ch07/instructions/base"
	"jvmgo/ch07/rtda"
	"jvmgo/ch07/rtda/heap"
)

// 取出类的某个静态变量值，然后推入栈顶
type GET_STATIC struct {
	// 常量池索引
	base.Index16Instruction
}

func (self *GET_STATIC) Execute(frame *rtda.Frame) {
	// 获取当前方法
	method := frame.Method()
	// 获取当前类
	currentClass := method.Class()
	// 获取常量池
	cp := currentClass.ConstantPool()
	// 获取字段引用
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
	// 解析字段引用，获取字段
	field := fieldRef.ResolvedField()
	// 获取这个字段属于哪个类
	class := field.Class()

	// 如果解析后的字段不是静态字段，要抛出IncompatibleClassChangeError异常。如果声明字段的类还没有初始化好，也需要先初始化。getstatic只是读取静态变量的值，自然也就不用管它是否是final了
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassError")
	}

	/*** 根据字段类型，从静态变量中取出相应的值，然后推入操作数栈顶  */
	// 获取字段描述符
	descriptor := field.Descriptor()
	// 获取这个字段存储在常量数组的索引值
	slotId := field.SlotId()
	// 类的常量数组
	slots := class.StaticVars()
	// 获取操作数栈
	stack := frame.OperandStack()
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I': //  Z--boolean, B--byte, C--Char, S--short, I--int
		stack.PushInt(slots.GetInt(slotId))
	case 'F': //  F--float
		stack.PushFloat(slots.GetFloat(slotId))
	case 'J': // J--long
		stack.PushLong(slots.GetLong(slotId))
	case 'D': // D--double
		stack.PushDouble(slots.GetDouble(slotId))
	case 'L', '[': // L--对象类型, [--数组
		stack.PushRef(slots.GetRef(slotId))
	default:
		// todo
	}
}
