package references

import (
	"jvmgo/ch07/instructions/base"
	"jvmgo/ch07/rtda"
	"jvmgo/ch07/rtda/heap"
)

/**
invokespecial指令用来调用无须动态绑定的实例方法，包括构造函数、私有方法和通过super关键字调用的超类方法
因为对象是需要初始化的，所以每个类都至少有一个构造函数。
即使用户自己不定义，编译器也会自动生成一个默认构造函数。
在创建类实例时，编译器会在 new指令 的后面加入 invokespecial 指令来调用构造函数初始化对象
*/
type INVOKE_SPECIAL struct{ base.Index16Instruction }

// 先拿到当前类、当前常量池、方法符号引用，然后解析符号引用，拿到解析后的类和方法
func (self *INVOKE_SPECIAL) Execute(frame *rtda.Frame) {
	// 获取当前类
	currentClass := frame.Method().Class()
	// 获取当前常量池
	cp := currentClass.ConstantPool()
	// 根据常量池获取方法符号引用
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	// 解析方法符号引用获取类和方法
	resolvedClass := methodRef.ResolvedClass()
	resolvedMethod := methodRef.ResolvedMethod()

	// 如果是构造方法, 则声明该方法的类必须是当前方法,即构造函数所属的类和解析后的类是一致的
	if resolvedMethod.Name() == "<init>" && resolvedMethod.Class() != resolvedClass {
		panic("java.lang.NoSuchMethodError")
	}

	// 如果是静态方法，则抛出IncompatibleClassChangeError异常，因为构造方法不是静态方法
	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 从操作数栈中弹出 this 引用
	ref := frame.OperandStack().GetRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if ref == nil {
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

	// 如果调用的是超类中的函数，但不是构造函数，且当前类的ACC_SUPER标志被设置，需要一个额外的过程查找最终要调用的方法；否则前面从方法符号引用中解析出来的方法就是要调用的方法。
	methodToBeInvoked := resolvedMethod
	if currentClass.IsSuper() &&
		resolvedClass.IsSuperClassOf(currentClass) &&
		resolvedMethod.Name() != "<init>" {
		// 需要一个额外的过程查找最终要调用的方法
		methodToBeInvoked = heap.LookupMethodInClass(currentClass.SuperClass(),
			methodRef.Name(), methodRef.Descriptor())
	}

	// 如果查找过程失败，或者找到的方法是抽象的，抛出AbstractMethodError异常
	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}

	base.InvokeMethod(frame, methodToBeInvoked)
}
