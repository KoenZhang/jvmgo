package heap

import (
	"jvmgo/ch06/classfile"
	"strings"
)

// name, superClassName and interfaceNames are all binary names(jvms8-4.2.1)
type Class struct {
	accessFlags       uint16        // 访问标志，一共16bit,具体含义在 access_flags.go
	name              string        // 类名称，形式类似于 java/lang/Object
	superClassName    string        // 超类名称，形式类似于 java/lang/Object
	interfaceNames    []string      // 实现的接口名称 ，形式类似于 java/lang/Object
	constantPool      *ConstantPool // 常量池，存放运行时常量池指针
	fields            []*Field      // 字段表
	methods           []*Method     // 方法表
	loader            *ClassLoader  // 存放类加载器指针
	superClass        *Class        // 超类
	interfaces        []*Class      // 实现的接口列表
	instanceSlotCount uint          // 实例变量占用的空间
	staticSlotCount   uint          // 类变量占用的空间
	staticVars        Slots         // 静态变量
}

/**
 * 把ClassFile结构体转换成Class结构体
 */
func newClass(cf *classfile.ClassFile) *Class {
	class := &Class{}
	class.accessFlags = cf.AccessFlags()
	class.name = cf.ClassName()
	class.superClassName = cf.SuperClassName()
	class.interfaceNames = cf.InterfaceNames()
	class.constantPool = newConstantPool(class, cf.ConstantPool())
	class.fields = newFields(class, cf.Fields())
	class.methods = newMethods(class, cf.Methods())
	return class
}

/**
getters
*/
func (self *Class) ConstantPool() *ConstantPool {
	return self.constantPool
}
func (self *Class) StaticVars() Slots {
	return self.staticVars
}

/**
以下8哥方法，均为判断访问标识是否被设置
*/

func (self *Class) IsPublic() bool {
	return 0 != self.accessFlags&ACC_PUBLIC
}
func (self *Class) IsFinal() bool {
	return 0 != self.accessFlags&ACC_FINAL
}
func (self *Class) IsSuper() bool {
	return 0 != self.accessFlags&ACC_SUPER
}
func (self *Class) IsInterface() bool {
	return 0 != self.accessFlags&ACC_INTERFACE
}
func (self *Class) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Class) IsSynthetic() bool {
	return 0 != self.accessFlags&ACC_SYNTHETIC
}
func (self *Class) IsAnnotation() bool {
	return 0 != self.accessFlags&ACC_ANNOTATION
}
func (self *Class) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}

// 如果类D想访问类C，需要满足两个条件之一：C是public，或者C和D在同一个运行时包内
func (self *Class) isAccessibleTo(other *Class) bool {
	return self.IsPublic() || self.getPackageName() == other.getPackageName()
}

// 如果类D想访问类C，需要满足两个条件之一：C是public，或者C和D在同一个运行时包内
func (self *Class) getPackageName() string {
	if i := strings.LastIndex(self.name, "/"); i > 0 {
		return self.name[:i]
	}
	return ""
}
