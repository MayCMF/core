package util

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

// MD5Hash - MD5 hash value
func MD5Hash(b []byte) string {
	h := md5.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5HashString - MD5 hash value
func MD5HashString(s string) string {
	return MD5Hash([]byte(s))
}

// SHA1Hash - SHA1 hash value
func SHA1Hash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1HashString - SHA1 hash value
func SHA1HashString(s string) string {
	return SHA1Hash([]byte(s))
}
