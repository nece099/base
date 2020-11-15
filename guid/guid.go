package guid

import (
	"fmt"
	"sync/atomic"
)

var gid = int64(10000)

func GenerateUID() int64 {
	newid := atomic.AddInt64(&gid, 1)
	return newid
}

func GenerateStringUID() string {
	s := fmt.Sprintf("%v", GenerateUID())
	return s
}
