package peer

import (
	"math/rand/v2"
	"time"

	"0xlayer/go-layer-mesh/utils/gopool"
)

type Request struct {
	date time.Time
}

func (p *Peer) IsRequest(id uint64) bool {
	return p.state.request.Has(id)
}

func (p *Peer) RequestId(id uint64) uint64 {
	p.state.request.Store(id, &Request{
		date: time.Now(),
	})
	return id
}

func (p *Peer) GetRequestId() uint64 {
	id := rand.Uint64()
	p.state.request.Store(id, &Request{
		date: time.Now(),
	})
	return id
}

func (p *Peer) RequestLoop() {
	gopool.Submit(func() {
		interval := (1 * time.Minute)
		timer := time.NewTimer(interval)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				var keys []uint64
				p.state.request.Range(func(k uint64, v *Request) bool {
					if time.Since(v.date) >= interval {
						keys = append(keys, k)
					}
					return true
				})
				for _, k := range keys {
					p.state.request.Delete(k)
				}
				timer.Reset(interval)
			case <-p.c:
				return
			}
		}
	})
}
