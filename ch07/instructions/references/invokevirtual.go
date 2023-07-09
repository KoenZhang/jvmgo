package references

import (
	"fmt"
	"jvmgo/ch07/instructions/base"
)
import "jvmgo/ch07/rtda"
import "jvmgo/ch07/rtda/heap"

/**
调用需要动态绑定的方法，针对非接口类型的引用调用方法，就使用invokeivirtual指令
使用super关键字调用超类中的方法不能使用invokevirtual
 当Java虚拟机通过invokevirtual调用方法时，this引用指向某个类（或其子类）的实例。
因为类的继承层次是固定的，所以虚拟机可以使用一种叫作vtable（Virtual Method Table）的技术加速方法查找
 Invoke instance method; dispatch based on class
*/
type INVOKE_VIRTUAL struct{ base.Index16Instruction }

func (self *INVOKE_VIRTUAL) Execute(frame *rtda.Frame) {
	// 获取当前类
	currentClass := frame.Method().Class()
	// 获取当前常量池
	cp := currentClass.ConstantPool()
	// 根据常量池获取方法符号引用
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	// 解析方法符号引用获取方法
	resolvedMethod := methodRef.ResolvedMethod()
	// 果是静态方法，则抛出IncompatibleClassChangeError异常，因为构造方法不是静态方法
	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 从操作数栈中弹出 this 引用
	ref := frame.OperandStack().GetRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if ref == nil {
		// hack!
		if methodRef.Name() == "println" {
			_println(frame.OperandStack(), methodRef.Descriptor())
			return
		}

		panic("java.lang.NullPointerException")
	}
	// protected 方法只能被声明该类方法的类或者子类调用
	if resolvedMethod.IsProtected() &&
		resolvedMethod.Class().IsSuperClassOf(currentClass) &&
		resolvedMethod.Class().GetPackageName() != currentClass.GetPackageName() &&
		ref.Class() != currentClass &&
		!ref.Class().IsSubClassOf(currentClass) {
		panic("java.lang.IllegalAccessError")
	}
	// 从对象的类中查找真正要调用的方法
	methodToBeInvoked := heap.LookupMethodInClass(ref.Class(),
		methodRef.Name(), methodRef.Descriptor())
	// 如果找不到方法，或者找到的是抽象方法，则需要抛出AbstractMethodError异常
	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}
	base.InvokeMethod(frame, methodToBeInvoked)
}

// hack!
func _println(stack *rtda.OperandStack, descriptor string) {
	switch descriptor {
	case "(Z)V":
		fmt.Printf("%v\n", stack.PopInt() != 0)
	case "(C)V":
		fmt.Printf("%c\n", stack.PopInt())
	case "(I)V", "(B)V", "(S)V":
		fmt.Printf("%v\n", stack.PopInt())
	case "(F)V":
		fmt.Printf("%v\n", stack.PopFloat())
	case "(J)V":
		fmt.Printf("%v\n", stack.PopLong())
	case "(D)V":
		fmt.Printf("%v\n", stack.PopDouble())
	default:
		panic("println: " + descriptor)
	}
	stack.PopRef()
}
