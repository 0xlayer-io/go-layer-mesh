package peer

import (
	"0xlayer/go-layer-mesh/utils/gopool"

	"github.com/ethereum/go-ethereum/common"
)

type Receiver struct {
	p    *Peer
	txs  chan common.Hash
	reqs chan func() error
	c    chan struct{}
}

const (
	maxTxsPerReq   = 256
	maxParallelReq = 10
)

func (r *Receiver) Transaction() {
	var (
		queue []common.Hash
		done  chan struct{}
	)

	for {
		qSize := len(queue)
		if done == nil && qSize > 0 {
			var (
				count int
				reqs  []common.Hash
				size  int
			)

			for count = 0; (count < qSize) && (size < maxTxsPerReq); count++ {
				hash := queue[count]
				reqs = append(reqs, hash)
				size++
			}

			queue = queue[count:]
			if size > 0 {
				done = make(chan struct{})
				gopool.Submit(func() {
					defer close(done)
					if err := r.p.SendRequestTxs(reqs); err != nil {
						// r.Close()
						return
					}
				})
			}
		}

		select {
		case <-done:
			done = nil
		case hash := <-r.txs:
			queue = append(queue, hash)
		case <-r.c:
			return
		}
	}
}

func (r *Receiver) Requests() {
	var (
		queue []func() error
		size  int
		done  = make(chan struct{}, maxParallelReq)
	)

	for {
		for len(queue) > 0 && size < maxParallelReq {
			req := queue[0]
			queue = queue[1:]
			size++

			gopool.Submit(func() {
				defer func() {
					done <- struct{}{}
				}()

				if err := req(); err != nil {
					// r.Close()
					return
				}
			})
		}

		select {
		case <-done:
			size--
		case req := <-r.reqs:
			queue = append(queue, req)
		case <-r.c:
			for i := 0; i < size; i++ {
				<-done
			}
			return
		}
	}
}

func (r *Receiver) Request() {
	var (
		queue []func() error
		done  chan struct{}
	)

	for {
		if done == nil && len(queue) > 0 {
			req := queue[0]
			queue = queue[1:]
			done = make(chan struct{})
			gopool.Submit(func() {
				defer close(done)
				if err := req(); err != nil {
					// r.Close()
					return
				}
			})
		}

		select {
		case <-done:
			done = nil
		case req := <-r.reqs:
			queue = append(queue, req)
		case <-r.c:
			return
		}
	}
}

func (r *Receiver) QueueTx(hash common.Hash) {
	if !r.Closed() {
		r.txs <- hash
	}
}

func (r *Receiver) QueueReq(req func() error) {
	if !r.Closed() {
		r.reqs <- req
	}
}

func (r *Receiver) Closed() bool {
	select {
	case <-r.c:
		return true
	default:
		return false
	}
}

func (r *Receiver) Close() {
	if !r.Closed() {
		close(r.c)
	}
}

func NewReceiver(p *Peer) *Receiver {
	r := &Receiver{
		p:    p,
		txs:  make(chan common.Hash, 1024),
		reqs: make(chan func() error, 128),
		c:    make(chan struct{}),
	}
	if p.BroadcastTx() {
		gopool.Submit(r.Transaction)
	}
	gopool.Submit(r.Request)
	return r
}
