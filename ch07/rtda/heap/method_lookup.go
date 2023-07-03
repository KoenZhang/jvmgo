package heap

// 在类的继承关系中寻找方法
func LookupMethodInClass(class *Class, name string, descriptor string) *Method {
	// 按照继承关系查找方法
	for c := class; c != nil; c = c.superClass {
		// 获取一个类的所有方法
		for _, method := range c.methods {
			if method.name == name && method.descriptor == descriptor {
				return method
			}
		}
	}

	return nil
}

// 递归查找：在类的接口中寻找方法，包括接口的多继承关系
func lookupMethodInInterfaces(ifaces []*Class, name string, descriptor string) *Method {
	// 遍历所有的接口
	for _, iface := range ifaces {
		// 获取每种接口的方法
		for _, method := range iface.methods {
			if method.name == name && method.descriptor == descriptor {
				return method
			}
		}

		// 接口的方法找不到，递归查找该接口继承的多个接口
		method := lookupMethodInInterfaces(iface.interfaces, name, descriptor)
		if method != nil {
			return method
		}
	}

	return nil
}
