package utils

import (
	"crypto/md5"
	"fmt"

	"github.com/OneOfOne/xxhash"

	"github.com/wrfly/yasuser/types"
)

func HashURL(url string, opts *types.ShortOptions) string {
	// custom URL connot exist with TTL
	in := fmt.Sprintf("%s:%s", url, opts.Pass)
	if opts.TTL.Seconds() != 0 {
		in = fmt.Sprintf("%s:%s:%s",
			url, opts.Pass, opts.TTL.String())
	}
	if opts.Custom != "" {
		in = fmt.Sprintf("%s:%s:%s",
			url, opts.Custom, opts.Pass)
	}
	return XXHash(in)
}

func XXHash(in string) string {
	return fmt.Sprintf("%x", xxhash.ChecksumString64(in))
}

func MD5(in string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(in)))
}
