package heap

import (
	"jvmgo/ch07/classfile"
)

type MethodRef struct {
	MemberRef
	method *Method
}

func newMethodRef(cp *ConstantPool, refInfo *classfile.ConstantMethodrefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

func (self *MethodRef) ResolvedMethod() *Method {
	if self.method == nil {
		self.resolveMethodRef()
	}
	return self.method
}

/**
如果类d想通过方法符号引用访问类C的某个方法：
1. 先要解析符号引用得到类C。
2. 如果C是接口，则抛出IncompatibleClassChangeError异常，
3. 根据方法名和描述符查找方法。如果找不到对应的方法，则抛出NoSuchMethodError异常，
4. 检查类D是否有权限访问该方法。如果没有，则抛出IllegalAccessError异常
*/
// jvms8 5.4.3.3
func (self *MethodRef) resolveMethodRef() {
	// 获取当前类
	d := self.cp.class
	// 1. 解析符号引用得到类C
	c := self.ResolvedClass()
	// 2. 如果C是接口，则抛出IncompatibleClassChangeError异常
	if c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 3. 根据方法名和描述符查找方法。如果找不到对应的方法，则抛出NoSuchMethodError异常，
	method := lookupMethod(c, self.name, self.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}
	// 4. 检查类D是否有权限访问该方法。如果没有，则抛出IllegalAccessError异常
	if !method.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}
	self.method = method
}

/**
先从C的继承层次中找，如果找不到，就去C的接口中找
*/
func lookupMethod(class *Class, name string, descriptor string) *Method {
	// 先从C的继承层次找
	method := LookupMethodInClass(class, name, descriptor)
	if method == nil {
		// 再从C的接口中找
		method = lookupMethodInInterfaces(class.interfaces, name, descriptor)
	}
	return method
}
