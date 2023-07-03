package heap

import "jvmgo/ch07/classfile"

type InterfaceMethodRef struct {
	MemberRef
	method *Method
}

func newInterfaceMethodRef(cp *ConstantPool, refInfo *classfile.ConstantInterfaceMethodrefInfo) *InterfaceMethodRef {
	ref := &InterfaceMethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

func (self *InterfaceMethodRef) ResolvedInterfaceMethod() *Method {
	if self.method == nil {
		self.resolveInterfaceMethodRef()
	}
	return self.method
}

/**
如果类d想通过方法符号引用访问类C的某个方法：
1. 先要解析符号引用得到类C。
2. 如果C不是接口，则抛出IncompatibleClassChangeError异常，
3. 根据方法名和描述符查找方法。如果找不到对应的方法，则抛出NoSuchMethodError异常，
4. 检查类D是否有权限访问该方法。如果没有，则抛出IllegalAccessError异常
*/
// jvms8 5.4.3.4
func (self *InterfaceMethodRef) resolveInterfaceMethodRef() {
	// 获取当前类
	d := self.cp.class
	// 1. 解析符号引用得到类C
	c := self.class
	// 2. 如果C不是接口，则抛出IncompatibleClassChangeError异常，
	if !c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	// 3. 根据方法名和描述符查找方法。如果找不到对应的方法，则抛出NoSuchMethodError异常，
	method := lookupInInterfacesMethod(c, self.name, self.descriptor)
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
1. 如果能在接口中找到方法，就返回找到的方法，
2. 否则调用lookupMethodInInterfaces（）函数在超接口中寻找
*/
func lookupInInterfacesMethod(iface *Class, name string, descriptor string) *Method {
	// 1. 如果能在接口中找到方法，就返回找到的方法
	for _, method := range iface.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}
	// 2. 否则调用lookupMethodInInterfaces（）函数在超接口中寻找
	return lookupMethodInInterfaces(iface.interfaces, name, descriptor)
}
