package classpath

import "io/ioutil"
import "path/filepath"

type DirEntry struct {
	absDir string // 存放目录绝对路径
}

func newDirEntry(path string) *DirEntry {
	absDir, err := filepath.Abs(path) // 获取绝对路径
	if err != nil {
		panic(err)
	}
	return &DirEntry{absDir}
}

func (self *DirEntry) readClass(className string) ([]byte, Entry, error) {
	fileName := filepath.Join(self.absDir, className)
	data, err := ioutil.ReadFile(fileName)
	return data, self, err
}

func (self *DirEntry) String() string {
	return self.absDir
}
