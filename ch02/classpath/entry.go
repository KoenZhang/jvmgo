package classpath

import "os"
import "strings"

const pathListSeparator = string(os.PathListSeparator) // 存放路径分割符

type Entry interface {
	// 寻找和加载Class文件, 参数为class相对路径, 路径间用 '/' 分割, 如 className = java/lang/Object.class
	readClass(className string) ([]byte, Entry, error)

	// 类似 Java 中的 toString(), 返回变量的字符串表示
	String() string
}

/* 根据 参数path 创建不同类型的 Entry 实例,一共四种：
newCompositeEntry: 组合模式, 由多个 entry 组合而成, 是数组类型
newWildcardEntry: 通配符格式的类路径, 实际上也是 newCompositeEntry 类型
newZipEntry: zip 或者 jar 文件格式的类路径
newDirEntry: 目录形式的类路径
*/
func newEntry(path string) Entry {
	if strings.Contains(path, pathListSeparator) {
		return newCompositeEntry(path)
	}
	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}
	return newDirEntry(path)
}
