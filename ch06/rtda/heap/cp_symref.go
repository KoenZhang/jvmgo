package heap

// 常量池： 类、字段、方法、接口 等引用，提取公共结构体
type SymRef struct {
	cp        *ConstantPool // 符号引用所在的运行时常量池指针，可以通过符号引用访问到运行时常量池，进一步又可以访问到类数据
	className string        // 类的完全限定名
	class     *Class        // 解析后的类结构体指针
}
