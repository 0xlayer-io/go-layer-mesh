package peer

import (
	"0xlayer/go-layer-mesh/eth/structs"
	"0xlayer/go-layer-mesh/utils/gopool"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Sender struct {
	p      *Peer
	blocks chan *structs.BlockAnnounce
	txs    chan *types.Transaction
	c      chan struct{}
}

const (
	maxQueuedTxAnns = 1024
	maxTxPacketSize = 8192
)

func (s *Sender) Block() {
	var (
		queue []*structs.BlockAnnounce
		done  chan struct{}
	)

	for {
		if done == nil && len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]
			if !s.IsKnownBlock(block.Hash, block.Number) {
				done = make(chan struct{})
				gopool.Submit(func() {
					defer close(done)
					if err := s.p.SendAnnounceBlock(block.Hash, block.Number); err != nil {
						// s.Close()
						return
					}
				})
			}
		}

		select {
		case <-done:
			done = nil
		case block := <-s.blocks:
			queue = append(queue, block)
		case <-s.c:
			return
		}
	}
}

func (s *Sender) Transaction() {
	var (
		queue []*types.Transaction
		done  chan struct{}
	)

	var (
		interval = (1 * time.Millisecond)
		timer    = time.NewTimer(interval)
		reset    = false
	)
	if s.p.IsTrusted() {
		interval = (100 * time.Microsecond)
	}
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			qSize := len(queue)
			if done == nil && qSize > 0 {
				var (
					count   int
					pending []*types.Transaction
					size    int
				)

				for count = 0; (count < qSize) && (size < maxTxPacketSize); count++ {
					tx := queue[count]
					if !s.IsKnownTx(tx.Hash()) {
						pending = append(pending, tx)
						size += int(tx.Size())
					}
				}
				queue = queue[:copy(queue, queue[count:])]

				if len(pending) > 0 {
					done = make(chan struct{})
					gopool.Submit(func() {
						defer close(done)
						if err := s.SendAnnounceTxs(pending); err != nil {
							// s.Close()
							return
						}
					})
				}
			}
		case <-done:
			done = nil
			reset = (len(queue) > 0)
			if reset {
				timer.Reset(interval)
			}
		case hash := <-s.txs:
			queue = append(queue, hash)
			if len(queue) > maxQueuedTxAnns {
				queue = queue[:copy(queue, queue[len(queue)-maxQueuedTxAnns:])]
			}
			if !reset {
				timer.Reset(interval)
				reset = true
			}
		case <-s.c:
			return
		}
	}
}

func (s *Sender) IsKnownBlock(hash common.Hash, number uint64) bool {
	if !s.p.IsKnownBlock(hash) && s.p.IsHeightBlock(number) {
		s.p.KnownBlock(hash, number)
		return false
	}
	return true
}

func (s *Sender) QueueBlock(block *structs.BlockAnnounce) {
	if !s.Closed() {
		s.blocks <- block
	}
}

func (s *Sender) QueueTx(tx *types.Transaction) {
	if !s.Closed() {
		s.txs <- tx
	}
}

func (s *Sender) SendAnnounceTxs(txs []*types.Transaction) error {
	p := s.p
	switch p.Version() {
	case 68:
		return p.SendAnnounceTxs68(txs)
	default:
		return p.SendAnnounceTxs66(txs)
	}
}

func (s *Sender) IsKnownTx(hash common.Hash) bool {
	if !s.p.IsKnownTx(hash) {
		s.p.KnownTx(hash)
		return false
	}
	return true
}

func (s *Sender) Closed() bool {
	select {
	case <-s.c:
		return true
	default:
		return false
	}
}

func (s *Sender) Close() {
	if !s.Closed() {
		close(s.c)
	}
}

func NewSender(p *Peer) *Sender {
	s := &Sender{
		p:      p,
		blocks: make(chan *structs.BlockAnnounce, 20),
		txs:    make(chan *types.Transaction, 1024),
		c:      make(chan struct{}),
	}
	gopool.Submit(s.Block)
	if p.BroadcastTx() {
		gopool.Submit(s.Transaction)
	}
	return s
}
