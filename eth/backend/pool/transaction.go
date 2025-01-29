package pool

import (
	"sync"
	"time"

	"0xlayer/go-layer-mesh/utils/gopool"
	"0xlayer/go-layer-mesh/utils/syncmap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxPool struct {
	tx   *types.Transaction
	send time.Time
	recv time.Time
	mu   sync.RWMutex
}

var txPools = syncmap.NewTypedMapOf[common.Hash, *TxPool](syncmap.CommonHasher)

func (t *TxPool) Tx() *types.Transaction {
	return t.tx
}

func (t *TxPool) Done() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tx != nil
}

func HasTx(hash common.Hash) bool {
	if t := GetTx(hash); t != nil {
		return t.Done()
	}
	return false
}

func GetTx(hash common.Hash) *TxPool {
	if t, ok := txPools.Load(hash); ok {
		return t
	}
	return nil
}

func AddTx(hash common.Hash, tx *types.Transaction) bool {
	if !txPools.Has(hash) {
		txPools.Store(hash, &TxPool{
			tx:   tx,
			recv: time.Now(),
		})
		return true
	}
	return false
}

func AddReqTx(hash common.Hash) bool {
	t := GetTx(hash)
	if t == nil {
		txPools.Store(hash, &TxPool{
			tx:   nil,
			send: time.Now(),
		})
		return true
	} else if !t.Done() {
		return true
	}
	return false
}

func UpdateTx(hash common.Hash, tx *types.Transaction) bool {
	if t := GetTx(hash); t != nil {
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.tx == nil {
			t.tx = tx
			t.recv = time.Now()
			return true
		}
	}
	return false
}

func TxLoop() {
	gopool.Submit(func() {
		interval := (5 * time.Minute)
		timer := time.NewTimer(interval)
		defer timer.Stop()

		for range timer.C {
			var keys []common.Hash
			txPools.Range(func(k common.Hash, v *TxPool) bool {
				if time.Since(v.recv) >= interval {
					keys = append(keys, k)
				}
				return true
			})
			for _, key := range keys {
				txPools.Delete(key)
			}
			timer.Reset(interval)
		}
	})
}
