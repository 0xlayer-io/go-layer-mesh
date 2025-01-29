package hashmap

import (
	"time"

	"0xlayer/go-layer-mesh/utils/syncmap"

	"github.com/ethereum/go-ethereum/common"
)

type HashesMap struct {
	hashs *syncmap.SyncMap[common.Hash, time.Time]
}

func New() *HashesMap {
	return &HashesMap{
		hashs: syncmap.NewTypedMapOf[common.Hash, time.Time](syncmap.CommonHasher),
	}
}

func (h *HashesMap) Size() int {
	return h.hashs.Size()
}

func (h *HashesMap) Contains(hash common.Hash) bool {
	return h.hashs.Has(hash)
}

func (h *HashesMap) Add(hash common.Hash) {
	h.hashs.Store(hash, time.Now())
}

func (h *HashesMap) Remove(hash common.Hash) {
	h.hashs.Delete(hash)
}

func (h *HashesMap) Expired(duration time.Duration) {
	if h.Size() > 0 {
		var keys []common.Hash
		h.hashs.Range(func(k common.Hash, v time.Time) bool {
			if time.Since(v) >= duration {
				keys = append(keys, k)
			}
			return true
		})
		for _, k := range keys {
			h.hashs.Delete(k)
		}
	}
}

func (h *HashesMap) Clear() {
	h.hashs.Clear()
}
