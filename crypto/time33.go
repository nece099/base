package crypto

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
)

func Time33(str string) int64 {
	hash := int64(5381) // 001 010 100 000 101 ,hash后的分布更好一些
	s := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	for i := 0; i < len(s); i++ {
		hash = hash*33 + int64(s[i])
	}
	return hash
}

func CalculateSig2(name string) uint64 {
	return uint64(Time33(name))
}

func CalculateSig(name string) uint32 {
	return crc32.ChecksumIEEE([]byte(name))
}

// func CalculateUUID() uint64 {
// 	uid := uuid.Must(uuid.NewV4()).String()
// 	return uint64(Time33(uid))
// }
