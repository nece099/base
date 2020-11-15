package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ChangeExt(fileName string, newExt string) string {
	// dir := filepath.Dir(path)
	base := filepath.Base(fileName)
	ext := filepath.Ext(base)
	filename := strings.TrimSuffix(base, ext)
	return filename + newExt
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		if os.IsNotExist(err) {
			return false
		}

		return false
	}

	return true
}

func MkdirIfNotExists(path string) error {
	if Exists(path) {
		return nil
	} else {
		return os.MkdirAll(path, os.ModePerm)
	}
}

func GetCurrentDirectory() string {

	// str, _ := os.Getwd()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
