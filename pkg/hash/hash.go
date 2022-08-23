package hash

import (
	"hash"
	"hash/fnv"
)

func Hash32(s string) hash.Hash32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h
}
