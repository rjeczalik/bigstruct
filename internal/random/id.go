package random

import (
	"hash/maphash"

	"github.com/google/uuid"
)

var (
	seed = maphash.MakeSeed()
)

func ID() uint64 {
	var (
		uuid = uuid.New()
		h    maphash.Hash
	)

	h.SetSeed(seed)
	h.WriteString(uuid.String())

	// Strip first high bit to workaround the database/sql driver error:
	//
	//   sql: uint64 values with high bit set are not supported
	//
	// More on this: https://github.com/golang/go/issues/6113
	//
	// return h.Sum64() >> 1
	return ^uint64(1<<63) & h.Sum64() // fixme: figure out which gives better distribution
}
