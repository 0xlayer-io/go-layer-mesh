package backend

import (
	"0xlayer/go-layer-mesh/eth/backend/pool"
	"0xlayer/go-layer-mesh/eth/handler/peer"
	"0xlayer/go-layer-mesh/eth/structs"
	"0xlayer/go-layer-mesh/utils/gopool"

	"github.com/ethereum/go-ethereum/common"
)

var (
	blockAnnounce = make(chan structs.BlockAnnounce, 20)
	txAnnounce    = make(chan common.Hash, 1024)
)

func AnnounceLoop() {
	gopool.Submit(func() {
		var (
			queue []structs.BlockAnnounce
			done  chan struct{}
		)

		for {
			if done == nil && len(queue) > 0 {
				block := queue[0]
				queue = queue[1:]
				done = make(chan struct{})
				gopool.Submit(func() {
					defer close(done)
					for _, p := range peer.Gets() {
						p.Sender().QueueBlock(&block)
					}
				})
			}

			select {
			case <-done:
				done = nil
			case number := <-blockAnnounce:
				queue = append(queue, number)
			}
		}
	})

	gopool.Submit(func() {
		var (
			queue []common.Hash
			done  chan struct{}
		)

		for {
			if done == nil && len(queue) > 0 {
				hash := queue[0]
				if tx := pool.GetTx(hash); tx != nil && tx.Tx() != nil {
					queue = queue[1:]
					done = make(chan struct{})
					gopool.Submit(func() {
						defer close(done)
						for _, p := range peer.Gets() {
							if p.BroadcastTx() {
								p.Sender().QueueTx(tx.Tx())
							}
						}
					})
				}
			}

			select {
			case <-done:
				done = nil
			case hash := <-txAnnounce:
				queue = append(queue, hash)
			}
		}
	})
}

func BlockAnnounce(hash common.Hash, number uint64) {
	blockAnnounce <- structs.BlockAnnounce{
		Hash:   hash,
		Number: number,
	}
}

func TxAnnounce(hash common.Hash) {
	txAnnounce <- hash
}
