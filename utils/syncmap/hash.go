package syncmap

import (
	"github.com/cespare/xxhash/v2"
	"github.com/ethereum/go-ethereum/common"

	"0xlayer/go-layer-mesh/p2p/enode"
)

type Hasher[K comparable] func(K, uint64) uint64

func EnodeHasher(k enode.ID, s uint64) uint64 {
	return writeBytesHash(k.Bytes(), s)
}

func CommonHasher(k common.Hash, s uint64) uint64 {
	return writeBytesHash(k.Bytes(), s)
}

func writeBytesHash(b []byte, s uint64) uint64 {
	h := xxhash.NewWithSeed(s)
	if _, err := h.Write(b); err != nil {
		panic(err)
	}
	return h.Sum64()
}
