package utils

import "github.com/nece099/base/crypto"

func HashPassword(pass string) string {
	return crypto.MD5("hash" + pass)
}
