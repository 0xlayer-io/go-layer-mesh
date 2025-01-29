package handler

import (
	handle "0xlayer/go-layer-mesh/eth/handler/message"
	"0xlayer/go-layer-mesh/eth/handler/peer"
	"0xlayer/go-layer-mesh/eth/message"
	"0xlayer/go-layer-mesh/p2p"
)

type handler func(*peer.Peer, p2p.Msg, uint32) error

var messages = map[uint64]handler{
	message.GetNodeDataMsg:                handle.GetNodeDataMsg,
	message.GetReceiptsMsg:                handle.GetReceiptsMsg,
	message.NewBlockMsg:                   handle.NewBlockMsg,
	message.NewBlockHashesMsg:             handle.NewBlockHashesMsg,
	message.GetBlockHeadersMsg:            handle.GetBlockHeadersMsg,
	message.GetBlockBodiesMsg:             handle.GetBlockBodiesMsg,
	message.BlockHeadersMsg:               handle.BlockHeadersMsg,
	message.BlockBodiesMsg:                handle.BlockBodiesMsg,
	message.TransactionsMsg:               handle.TransactionsMsg,
	message.NewPooledTransactionHashesMsg: handle.NewPooledTransactionHashesMsg,
	message.GetPooledTransactionsMsg:      handle.GetPooledTransactionsMsg,
	message.PooledTransactionsMsg:         handle.PooledTransactionsMsg,
}

func Message(p *peer.Peer, msg p2p.Msg, version uint32) (bool, error) {
	if version == 66 || version == 67 || version == 68 {
		if h, ok := messages[msg.Code]; ok {
			return true, h(p, msg, version)
		}
	}
	return false, nil
}
