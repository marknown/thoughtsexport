package utils

import (
	"io/ioutil"
	"os"
)

// FileRead 读取文件内容
func FileRead(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if nil != err {
		return ""
	}

	return string(bytes)
}

// FileAppend 把字符串追加到指定文件
func FileAppend(fileName string, s string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	buf := []byte(s)
	fd.Write(buf)
	fd.Close()
}

// FileWrite 把字符串追写到指定文件，之前的内容会被清空
func FileWrite(fileName string, s string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	buf := []byte(s)
	fd.Write(buf)
	fd.Close()
}

// FileExist 判断路径或者文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
