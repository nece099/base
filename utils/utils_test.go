package utils

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestUtils(t *testing.T) {
	str := ChangeExt("v111222.mp4", ".webm")
	fmt.Print(str)
}

func TestUtils2(t *testing.T) {
	path := "/Users/dev/workspace/golang/gopath/src/"

	dir := filepath.Dir(path)

	fmt.Print(dir)
}

func TestPassword(t *testing.T) {
	pass := "123456" //0562b36c3c5a3925dbe3c4d32a4f2ba2
	hash := HashPassword(pass)
	fmt.Print(hash)
}
