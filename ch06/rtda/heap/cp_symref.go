package heap

// 常量池： 类、字段、方法、接口 等引用，提取公共结构体
type SymRef struct {
	cp        *ConstantPool // 符号引用所在的运行时常量池指针，可以通过符号引用访问到运行时常量池，进一步又可以访问到类数据
	className string        // 类的完全限定名
	class     *Class        // 解析后的类结构体指针
}

// 如果类符号引用已经解析，ResolvedClass（）方法直接返回类指针
func (self *SymRef) ResolvedClass() *Class {
	if self.class == nil {
		self.ResolveClassRef()
	}
	return self.class
}

// 如果类D通过符号引用N引用类C的话，要解析N，先用D的类加载器加载C，然后检查D是否有权限访问C，如果没有，则抛出IllegalAccessError异常
func (self *SymRef) ResolveClassRef() {
	d := self.class
	c := d.loader.LoadClass(self.className)
	if !c.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}
	self.class = c
}
