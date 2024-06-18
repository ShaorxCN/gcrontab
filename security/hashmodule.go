package security

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

// HashSha256 返回16进制的sha256结果
func HashSha256(s string) string {
	res := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", res[:])
}

func HashMD5(s string) string {
	res := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", res[:])
}
