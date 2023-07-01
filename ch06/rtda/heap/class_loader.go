package heap

import (
	"fmt"
	"jvmgo/ch06/classfile"
	"jvmgo/ch06/classpath"
)

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
	fmt.Printf("[Loaded %s from %s]\n", name, entry)
	return class
}

// 类的链接分为验证和准备两个阶段
func link(class *Class) {
	verify(class)
	prepare(class)
}

func prepare(class *Class) {

}

// 为了确保安全性，Java虚拟机规范要求在执行类的任何代码之前，对类进行严格的验证
func verify(class *Class) {
	// todo
}

func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := self.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException:" + name)
	}
	return data, entry
}

func (self *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	class.loader = self
	resolveSuperClass(class)
	resolveInterfaces(class)
	self.classMap[class.name] = class
	return class
}

func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

func resolveSuperClass(class *Class) {
	// 所有的类都只有一个超类，Object没有父类
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		panic("java.lang.ClassFormatError")
	}
	return newClass(cf)
}
