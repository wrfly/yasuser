package utils

import (
	"crypto/md5"
	"fmt"

	"github.com/OneOfOne/xxhash"
)

func XXHash(in string) string {
	return fmt.Sprintf("%x", xxhash.ChecksumString64(in))
}

func MD5(in string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(in)))
}
