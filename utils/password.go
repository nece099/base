package utils

import "github.com/zen099/onetube/server/base/crypto"

func HashPassword(pass string) string {
	return crypto.MD5("hash" + pass)
}
