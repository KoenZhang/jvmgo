package references

import (
	"jvmgo/ch06/instructions/base"
	"jvmgo/ch06/rtda"
	"jvmgo/ch06/rtda/heap"
)

/**
putstatic指令给类的某个静态变量赋值，
它需要两个操作数:
1. 第一个操作数是uint16索引，来自字节码。通过这个索引可以从当前类的运行时常量池中找到一个字段符号引用，解析这个符号引用就可以知道要给类的哪个静态变量赋值。
2. 第二个操作数是要赋给静态变量的值，从操作数栈中弹出
*/

type PUT_STATIC struct {
	// 静态变量在运行时常量池的索引
	base.Index16Instruction
}

/**
putstatic指令给类的某个静态变量赋值，
它需要两个操作数:
1. 第一个操作数是uint16索引，来自字节码。通过这个索引可以从当前类的运行时常量池中找到一个字段符号引用，解析这个符号引用就可以知道要给类的哪个静态变量赋值。
2. 第二个操作数是要赋给静态变量的值，从操作数栈中弹出
*/

func (self *PUT_STATIC) Execute(frame *rtda.Frame) {
	// 拿到当前方法
	currentMethod := frame.Method()
	// 根据当前方法获取当前类
	currentClass := currentMethod.Class()
	// 根据当前类获取当前常量池
	cp := currentClass.ConstantPool()
	//  获取字段符号引用
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
	// 解析字段符号引用
	field := fieldRef.ResolvedField()
	// 获取当前字段所属的类型
	class := field.Class()
	// 非静态字段，抛出异常
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 常量， 只能在类初始化方法中给它赋值
	if field.IsFinal() {
		if currentClass != class || currentMethod.Name() != "<clinit>" {
			panic("java.lang.IllegalAccessError")
		}
	}

	// 根据字段类型从操作数栈中弹出相应的值，然后赋给静态变量
	descriptor := field.Descriptor()
	slotId := field.SlotId()
	slots := class.StaticVars()
	stack := frame.OperandStack()

	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I': // Z--boolean, B--byte, C--Char, S--short, I--int
		slots.SetInt(slotId, stack.PopInt())
	case 'F': // F--float
		slots.SetFloat(slotId, stack.PopFloat())
	case 'J': // J--long
		slots.SetLong(slotId, stack.PopLong())
	case 'D': // D--double
		slots.SetDouble(slotId, stack.PopDouble())
	case 'L', '[': // L--对象类型, [--数组
		slots.SetRef(slotId, stack.PopRef())
	default: // V--void
		// todo
	}
}
