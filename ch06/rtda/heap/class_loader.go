package heap

import "jvmgo/ch06/classpath"

type ClassLoader struct {
	cp       *classpath.Classpath // class文件所在路径
	classMap map[string]*Class    //	已经加载过的class文件， key是类的完全限定名， 这里可以将其当作方法区的简化实现
}

func newClassLoader(cp *classpath.Classpath) *ClassLoader {
	return &ClassLoader{
		cp:       cp,
		classMap: make(map[string]*Class),
	}
}

// 把类数据加载到方法区. 先查找classMap，看类是否已经被加载。如果是，直接返回类数据，否则调用loadNonArrayClass（）方法加载类
func (self *ClassLoader) LoadClass(name string) *Class {
	if class, ok := self.classMap[name]; ok {
		return class // 类已经加载
	}
	return self.loadNonArrayClass(name)
}

// 数组类和普通类有很大的不同，它的数据并不是来自class文件，而是由Java虚拟机在运行期间生成
func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry := self.readClass(name) // 读取文件
	class := self.defineClass(data)     // 解析class文件，生成虚拟机可以使用的类数据，并放入方法区
	link(class)                         // 链接
	return class
}
