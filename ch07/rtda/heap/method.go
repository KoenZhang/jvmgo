package heap

import "jvmgo/ch07/classfile"

type Method struct {
	ClassMember
	maxStack     uint   // 操作数栈大小
	maxLocals    uint   // 局部变量表大小
	code         []byte // 方法字节码
	argSlotCount uint   // 方法的参数个数
}

func (self *Method) SetArgSlotCount(argSlotCount uint) {
	self.argSlotCount = argSlotCount
}

func (self *Method) ArgSlotCount() uint {
	return self.argSlotCount
}

func (self *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil {
		self.maxStack = codeAttr.MaxStack()
		self.maxLocals = codeAttr.MaxLocals()
		self.code = codeAttr.Code()
	}
}

func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = &Method{}
		methods[i].class = class
		methods[i].copyMemberInfo(cfMethod)
		methods[i].copyAttributes(cfMethod)
		methods[i].calcuArgSlotCount()
	}
	return methods
}

func (self *Method) IsSynchronized() bool {
	return 0 != self.accessFlags&ACC_SYNCHRONIZED
}
func (self *Method) IsBridge() bool {
	return 0 != self.accessFlags&ACC_BRIDGE
}
func (self *Method) IsVarargs() bool {
	return 0 != self.accessFlags&ACC_VARARGS
}
func (self *Method) IsNative() bool {
	return 0 != self.accessFlags&ACC_NATIVE
}
func (self *Method) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Method) IsStrict() bool {
	return 0 != self.accessFlags&ACC_STRICT
}

// getters
func (self *Method) MaxStack() uint {
	return self.maxStack
}
func (self *Method) MaxLocals() uint {
	return self.maxLocals
}
func (self *Method) Code() []byte {
	return self.code
}

// 计算 ArgSlotCount
func (self *Method) calcuArgSlotCount() {
	// 分解方法描述符
	parsedDescriptor := parseMethodDescriptor(self.descriptor)
	for _, paramType := range parsedDescriptor.parameterTypes {
		self.argSlotCount++
		// long(J)和 double(D) 类型占用 2 个位置
		if paramType == "J" || paramType == "D" {
			self.argSlotCount++
		}
	}

	// 实例方法多增加一个参数 -> this
	if !self.IsStatic() {
		self.argSlotCount++
	}

}
