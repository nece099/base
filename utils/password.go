package utils

import "github.com/nece099/base/encrypt"

func HashPassword(pass string) string {
	return encrypt.MD5("hash" + pass)
}
