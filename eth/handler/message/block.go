package handler_message

import (
	"fmt"

	"0xlayer/go-layer-mesh/eth/backend"
	"0xlayer/go-layer-mesh/eth/backend/pool"
	"0xlayer/go-layer-mesh/eth/handler/peer"
	"0xlayer/go-layer-mesh/eth/message/packet"
	"0xlayer/go-layer-mesh/p2p"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func NewBlockMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.NewBlockPacket
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("NewBlockMsg: %v", err)
	}

	hash, number, header, body := decode.Block.Hash(), decode.Block.NumberU64(), decode.Block.Header(), &packet.BlockBody{
		Transactions: decode.Block.Transactions(),
		Uncles:       decode.Block.Uncles(),
		Withdrawals:  decode.Block.Withdrawals(),
		Sidecars:     decode.Sidecars,
	}

	p.KnownBlock(hash, number)
	if pool.AddBlock(hash, header, body) {
		backend.BlockAnnounce(hash, number)
		p.BlockSync()
	}
	return nil
}

func NewBlockHashesMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.NewBlockHashesPacket
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("NewBlockHashesMsg: %v", err)
	}

	for _, v := range decode {
		if v.Hash != (common.Hash{}) && v.Number > 0 {
			p.KnownBlock(v.Hash, v.Number)
			if !pool.IsBlock(v.Hash) {
				if block := pool.GetBlock(v.Hash); block != nil {
					p.SendRequestBlock(v.Number, []common.Hash{v.Hash}, !block.IsHeader(), !block.IsBody())
				} else if pool.AddReqBlock(v.Hash, v.Number) {
					p.SendRequestBlock(v.Number, []common.Hash{v.Hash}, true, true)
				}
			}
		}
	}
	return nil
}

func GetBlockHeadersMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.GetBlockHeadersPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("GetBlockHeadersMsg: %v", err)
	}

	request := decode.GetBlockHeadersPacket
	if request.Amount != 1 || request.Skip != 0 || request.Reverse {
		p.SendEmptyBlockHeaders(decode.RequestId)
		return nil
	}

	var block *pool.BlockPool
	if request.Origin.Number != 0 {
		block = pool.GetBlockNumber(request.Origin.Number)
	} else if request.Origin.Hash != (common.Hash{}) {
		block = pool.GetBlock(request.Origin.Hash)
	}

	if block != nil && block.IsHeader() {
		p.SendBlockHeaders([]*types.Header{block.Header()}, decode.RequestId)
	} else {
		p.SendEmptyBlockHeaders(decode.RequestId)
	}
	return nil
}

func GetBlockBodiesMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.GetBlockBodiesPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("GetBlockBodiesMsg: %v", err)
	}

	var bodys []*packet.BlockBody
	for _, hash := range decode.GetBlockBodiesPacket {
		if block := pool.GetBlock(hash); block != nil {
			if block.IsBody() {
				bodys = append(bodys, block.Body())
			}
		}
	}

	if len(bodys) > 0 {
		p.SendBlockBodies(bodys, decode.RequestId)
	} else {
		p.SendEmptyBlockBodies(decode.RequestId)
	}
	return nil
}

func BlockHeadersMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.BlockHeadersPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("BlockHeadersMsg: %v", err)
	}

	if p.IsRequest(decode.RequestId) {
		for _, v := range decode.BlockHeadersPacket {
			if v != nil {
				if pool.UpdateBlockHeader(v.Hash(), v) {
					backend.BlockAnnounce(v.Hash(), v.Number.Uint64())
					p.BlockSync()
				}
			}
		}
	}
	return nil
}

func BlockBodiesMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.BlockBodiesPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("BlockBodiesMsg: %v", err)
	}

	if p.IsRequest(decode.RequestId) {
		for _, v := range decode.BlockBodiesPacket {
			if v != nil {
				if s, h := pool.UpdateBlockBody(decode.RequestId, v); s {
					backend.BlockAnnounce(h, decode.RequestId)
					p.BlockSync()
				}
			}
		}
	}
	return nil
}
