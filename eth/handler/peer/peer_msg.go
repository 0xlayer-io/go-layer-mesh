package peer

import (
	"0xlayer/go-layer-mesh/eth/message"
	"0xlayer/go-layer-mesh/eth/message/packet"
	"0xlayer/go-layer-mesh/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func (p *Peer) SendEmptyNodeData(requestId uint64) {
	p.Receiver().QueueReq(func() error {
		var empty [][]byte
		return p.Send(message.NodeDataMsg, &packet.NodeDataPacket66{
			RequestId:      requestId,
			NodeDataPacket: packet.NodeDataPacket(empty),
		})
	})
}

func (p *Peer) SendEmptyReceipts(requestId uint64) {
	p.Receiver().QueueReq(func() error {
		var empty []rlp.RawValue
		return p.Send(message.ReceiptsMsg, &packet.ReceiptsRLPPacket66{
			RequestId:         requestId,
			ReceiptsRLPPacket: packet.ReceiptsRLPPacket(empty),
		})
	})
}

func (p *Peer) SendBlockHeaders(headers []*types.Header, requestId uint64) {
	p.Receiver().QueueReq(func() error {
		reqs := utils.EncodeToRLPs(headers)
		return p.Send(message.BlockHeadersMsg, &packet.BlockHeadersRLPPacket66{
			RequestId:             requestId,
			BlockHeadersRLPPacket: packet.BlockHeadersRLPPacket(reqs),
		})
	})
}

func (p *Peer) SendEmptyBlockHeaders(requestId uint64) {
	p.Receiver().QueueReq(func() error {
		var empty []rlp.RawValue
		return p.Send(message.BlockHeadersMsg, &packet.BlockHeadersRLPPacket66{
			RequestId:             requestId,
			BlockHeadersRLPPacket: packet.BlockHeadersRLPPacket(empty),
		})
	})
}

func (p *Peer) SendBlockBodies(bodys []*packet.BlockBody, requestId uint64) {
	p.Receiver().QueueReq(func() error {
		reqs := utils.EncodeToRLPs(bodys)
		return p.Send(message.BlockBodiesMsg, &packet.BlockBodiesRLPPacket66{
			RequestId:            requestId,
			BlockBodiesRLPPacket: packet.BlockBodiesRLPPacket(reqs),
		})
	})
}

func (p *Peer) SendEmptyBlockBodies(requestId uint64) {
	p.Receiver().QueueReq(func() error {
		var empty []rlp.RawValue
		return p.Send(message.BlockBodiesMsg, &packet.BlockBodiesRLPPacket66{
			RequestId:            requestId,
			BlockBodiesRLPPacket: packet.BlockBodiesRLPPacket(empty),
		})
	})
}

func (p *Peer) SendRequestBlockHeader(number uint64) {
	p.Receiver().QueueReq(func() error {
		return p.Send(message.GetBlockHeadersMsg, &packet.GetBlockHeadersPacket66{
			RequestId: p.GetRequestId(),
			GetBlockHeadersPacket: &packet.GetBlockHeadersPacket{
				Origin:  packet.HashOrNumber{Number: number},
				Amount:  1,
				Skip:    0,
				Reverse: false,
			},
		})
	})
}

func (p *Peer) SendRequestBlockBodys(number uint64, hash []common.Hash) {
	p.Receiver().QueueReq(func() error {
		return p.Send(message.GetBlockBodiesMsg, &packet.GetBlockBodiesPacket66{
			RequestId:            p.RequestId(number),
			GetBlockBodiesPacket: packet.GetBlockBodiesPacket(hash),
		})
	})
}

func (p *Peer) SendRequestBlock(number uint64, hash []common.Hash, header, body bool) {
	if header {
		p.SendRequestBlockHeader(number)
	}
	if body {
		p.SendRequestBlockBodys(number, hash)
	}
}

func (p *Peer) SendAnnounceBlock(hash common.Hash, number uint64) error {
	return p.Send(message.NewBlockHashesMsg, &packet.NewBlockHashesPacket{{hash, number}})
}

func (p *Peer) SendRequestTxs(txs []common.Hash) error {
	return p.Send(message.GetPooledTransactionsMsg, &packet.GetPooledTransactionsPacket66{
		RequestId:                   p.GetRequestId(),
		GetPooledTransactionsPacket: packet.GetPooledTransactionsPacket(txs),
	})
}

func (p *Peer) SendAnnounceTxs66(txs []*types.Transaction) error {
	var hashes []common.Hash
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash())
	}
	return p.Send(message.NewPooledTransactionHashesMsg, packet.NewPooledTransactionHashesPacket(hashes))
}

func (p *Peer) SendAnnounceTxs68(txs []*types.Transaction) error {
	var types []byte
	var sizes []uint32
	var hashes []common.Hash
	for _, tx := range txs {
		types = append(types, tx.Type())
		sizes = append(sizes, uint32(tx.Size()))
		hashes = append(hashes, tx.Hash())
	}
	return p.Send(message.NewPooledTransactionHashesMsg, &packet.NewPooledTransactionHashesPacket68{
		Types:  types,
		Sizes:  sizes,
		Hashes: hashes,
	})
}

func (p *Peer) SendPooledTxs(txs []*types.Transaction, requestId uint64) {
	if reqs := utils.EncodeToRLPs(txs); len(reqs) > 0 {
		p.Receiver().QueueReq(func() error {
			return p.Send(message.PooledTransactionsMsg, &packet.PooledTransactionsRLPPacket66{
				RequestId:                   requestId,
				PooledTransactionsRLPPacket: packet.PooledTransactionsRLPPacket(reqs),
			})
		})
	}
}
