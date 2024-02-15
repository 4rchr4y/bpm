package bundleutil

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

func ChecksumSHA256[T ~string | []byte](h hash.Hash, data T) string {
	switch v := any(data).(type) {
	case string:
		h.Write([]byte(v))

	case []byte:
		h.Write(v)
	}

	var buf [sha256.Size]byte
	sum := h.Sum(buf[:0])
	return hex.EncodeToString(sum)
}
