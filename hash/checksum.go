package hash

import (
	"bytes"
	"crypto/md5"
)

// Which returns MD5 checksum.
func ComputeMd5Checksum(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func ValidateMd5Checksum(data []byte, sum []byte) bool {
	ourSum := ComputeMd5Checksum(data)
	return bytes.Equal(ourSum, sum)
}
