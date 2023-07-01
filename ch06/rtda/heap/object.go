package heap

type Object struct {
	class  *Class // 对象的Class指针
	fields Slots  // 实例变量
}
