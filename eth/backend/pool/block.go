package pool

import (
	"sync"
	"time"

	"0xlayer/go-layer-mesh/eth/message/packet"
	"0xlayer/go-layer-mesh/utils/gopool"
	"0xlayer/go-layer-mesh/utils/syncmap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockPool struct {
	number uint64
	header *types.Header
	body   *packet.BlockBody
	recv   time.Time
	mu     sync.RWMutex
}

var blockPools = syncmap.NewTypedMapOf[common.Hash, *BlockPool](syncmap.CommonHasher)

func (b *BlockPool) Header() *types.Header {
	return b.header
}

func (b *BlockPool) Body() *packet.BlockBody {
	return b.body
}

func (b *BlockPool) IsHeader() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.header != nil
}

func (b *BlockPool) IsBody() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.body != nil
}

func IsBlock(hash common.Hash) bool {
	if b := GetBlock(hash); b != nil {
		b.mu.RLock()
		defer b.mu.RUnlock()
		return b.header != nil && b.body != nil
	}
	return false
}

func GetBlock(hash common.Hash) *BlockPool {
	if b, ok := blockPools.Load(hash); ok {
		return b
	}
	return nil
}

func GetBlockNumber(number uint64) *BlockPool {
	var block *BlockPool
	blockPools.Range(func(k common.Hash, v *BlockPool) bool {
		if v.number == number {
			block = v
			return false
		}
		return true
	})
	return block
}

func GetBlockRequest(number uint64) *BlockPool {
	var block *BlockPool
	blockPools.Range(func(k common.Hash, v *BlockPool) bool {
		if v.body == nil && v.number == number {
			block = v
			return false
		}
		return true
	})
	return block
}

func AddBlock(hash common.Hash, header *types.Header, body *packet.BlockBody) bool {
	if !blockPools.Has(hash) {
		blockPools.Store(hash, &BlockPool{
			header: header,
			body:   body,
			recv:   time.Now(),
		})
		return true
	}
	return false
}

func AddReqBlock(hash common.Hash, number uint64) bool {
	if !blockPools.Has(hash) {
		blockPools.Store(hash, &BlockPool{
			number: number,
			header: nil,
			body:   nil,
		})
		return true
	}
	return false
}

func UpdateBlockHeader(hash common.Hash, header *types.Header) bool {
	if b := GetBlock(hash); b != nil {
		b.mu.Lock()
		defer b.mu.Unlock()
		if b.header == nil {
			b.header = header
		}
		done := b.body != nil
		if done {
			b.recv = time.Now()
		}
		return done
	}
	return false
}

func UpdateBlockBody(number uint64, body *packet.BlockBody) (bool, common.Hash) {
	if b := GetBlockRequest(number); b != nil {
		b.mu.Lock()
		defer b.mu.Unlock()
		if b.body == nil {
			b.body = body
		}
		if b.header != nil {
			b.recv = time.Now()
			return true, b.header.Hash()
		}
	}
	return false, common.Hash{}
}

func BlockLoop() {
	gopool.Submit(func() {
		interval := (5 * time.Minute)
		timer := time.NewTimer(interval)
		defer timer.Stop()

		for range timer.C {
			var keys []common.Hash
			blockPools.Range(func(k common.Hash, v *BlockPool) bool {
				if time.Since(v.recv) >= interval {
					keys = append(keys, k)
				}
				return true
			})
			for _, k := range keys {
				blockPools.Delete(k)
			}
			timer.Reset(interval)
		}
	})
}
