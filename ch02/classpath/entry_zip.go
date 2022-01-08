package classpath

import "archive/zip"
import "errors"
import "io/ioutil"
import "path/filepath"

type ZipEntry struct {
	absPath string
}

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath}
}

func (self *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	r, err := zip.OpenReader(self.absPath) // 打开压缩文件
	if err != nil {                        // 打开压缩文件失败，直接返回
		return nil, nil, err
	}
	defer r.Close() // 确保打开的文件一定可以关闭

	// 遍历压缩文件中的文件，寻找 className
	for _, f := range r.File {
		if f.Name == className { // 找到对应的文件后，开始读取文件
			rc, err := f.Open()
			if err != nil {
				return nil, nil, err
			}
			defer rc.Close() // 确保打开的文件一定可以关闭

			data, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, nil, err
			}
			return data, self, nil
		}
	}

	return nil, nil, errors.New("Class not found: " + className)
}

func (self *ZipEntry) String() string {
	return self.absPath
}
