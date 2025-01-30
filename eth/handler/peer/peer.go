package peer

import (
	"errors"
	"sync/atomic"
	"time"

	"0xlayer/go-layer-mesh/p2p"
	"0xlayer/go-layer-mesh/p2p/enode"
	"0xlayer/go-layer-mesh/utils/gopool"
	"0xlayer/go-layer-mesh/utils/hashmap"
	"0xlayer/go-layer-mesh/utils/syncmap"

	"github.com/ethereum/go-ethereum/common"
)

type Peer struct {
	p       *p2p.Peer
	rw      p2p.MsgReadWriter
	info    *PeerInfo
	msg     *PeerMessage
	state   *PeerState
	stats   *PeerStats
	sync    *PeerSync
	ext     *PeerExtension
	trusted bool
	c       chan struct{}
}

type PeerInfo struct {
	Name    string
	Enode   string
	Version uint
}

type PeerMessage struct {
	sender   *Sender
	receiver *Receiver
}

type PeerState struct {
	txs         *hashmap.HashesMap
	blocks      *hashmap.HashesMap
	request     *syncmap.SyncMap[uint64, *Request]
	blockHeight atomic.Uint64
}

type PeerStats struct {
	SharedTxs      uint64
	AnnounceTxs    uint64
	SharedBlocks   uint64
	AnnounceBlocks uint64
	TotalRequest   uint64
	TotalResponse  uint64
}

type PeerSync struct {
	block       time.Time
	transaction time.Time
}

type PeerExtension struct {
	BroadcastTx bool
}

var peers = syncmap.NewTypedMapOf[enode.ID, *Peer](syncmap.EnodeHasher)

func (p *Peer) Peer() *p2p.Peer {
	return p.p
}

func (p *Peer) Msg() p2p.MsgReadWriter {
	return p.rw
}

func (p *Peer) Version() uint {
	return p.info.Version
}

func (p *Peer) Send(msgCode uint64, msg interface{}) error {
	if p.Closed() {
		return errors.New("peer is closed")
	}
	return p2p.Send(p.rw, msgCode, msg)
}

func (p *Peer) Receiver() *Receiver {
	return p.msg.receiver
}

func (p *Peer) Sender() *Sender {
	return p.msg.sender
}

func (p *Peer) IsTrusted() bool {
	return p.trusted
}

func (p *Peer) IsKnownTx(hash common.Hash) bool {
	return p.state.txs.Contains(hash)
}

func (p *Peer) KnownTx(hash common.Hash) {
	p.state.txs.Add(hash)
}

func (p *Peer) IsHeightBlock(number uint64) bool {
	return p.state.blockHeight.Load() <= number
}

func (p *Peer) IsKnownBlock(hash common.Hash) bool {
	return p.state.blocks.Contains(hash)
}

func (p *Peer) KnownBlock(hash common.Hash, number uint64) {
	if p.state.blockHeight.Load() < number {
		p.state.blockHeight.Store(number)
	}
	p.state.blocks.Add(hash)
}

func (p *Peer) BroadcastTx() bool {
	return p.ext.BroadcastTx
}

func (p *Peer) TxSync() {
	p.sync.transaction = time.Now()
}

func (p *Peer) BlockSync() {
	p.sync.block = time.Now()
}

func (p *Peer) Loop() {
	// block, transaction (expired)
	expire := (1 * time.Minute)
	gopool.Submit(func() {
		timer := time.NewTicker(expire)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				p.state.blocks.Expired(expire)
				p.state.txs.Expired(expire)
			case <-p.c:
				return
			}
		}
	})

	// useless peer (disconnect)
	if !p.IsTrusted() {
		gopool.Submit(func() {
			interval := (15 * time.Second)
			timer := time.NewTimer(interval)
			defer timer.Stop()
			for {
				select {
				case <-timer.C:
					tSync := time.Since(p.sync.transaction)
					if tSync > (10 * time.Minute) {
						p.p.Disconnect(p2p.DiscUselessPeer)
						return
					}
					timer.Reset(interval)
				case <-p.c:
					return
				}
			}
		})
	}
}

func (p *Peer) Open() {
	peers.Store(p.Peer().ID(), p)
	p.msg = &PeerMessage{
		sender:   NewSender(p),
		receiver: NewReceiver(p),
	}
	p.RequestLoop()
	p.Loop()
}

func (p *Peer) Close() {
	peers.Delete(p.Peer().ID())
	p.msg.sender.Close()
	p.msg.receiver.Close()
	p.state.request.Clear()
	p.state.blocks.Clear()
	p.state.txs.Clear()
	close(p.c)
}

func (p *Peer) Closed() bool {
	select {
	case <-p.c:
		return true
	default:
		return false
	}
}

func Has(id enode.ID) bool {
	return peers.Has(id)
}

func Get(id enode.ID) *Peer {
	if v, ok := peers.Load(id); ok {
		return v
	}
	return nil
}

func Gets() []*Peer {
	var result []*Peer
	peers.Range(func(k enode.ID, v *Peer) bool {
		result = append(result, v)
		return true
	})
	return result
}

func Size() uint32 {
	return uint32(peers.Size())
}

func New(peer *p2p.Peer, rw p2p.MsgReadWriter, trusted bool) *Peer {
	return &Peer{
		p:    peer,
		rw:   rw,
		info: &PeerInfo{},
		msg:  &PeerMessage{},
		state: &PeerState{
			txs:     hashmap.New(),
			blocks:  hashmap.New(),
			request: syncmap.NewIntegerMapOf[uint64, *Request](),
		},
		stats: &PeerStats{},
		sync: &PeerSync{
			block:       time.Now(),
			transaction: time.Now(),
		},
		ext: &PeerExtension{
			BroadcastTx: true,
		},
		trusted: trusted,
		c:       make(chan struct{}),
	}
}
