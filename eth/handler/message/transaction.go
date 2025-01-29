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

func TransactionsMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.TransactionsPacket
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("TransactionsMsg: %v", err)
	}

	var hashes []common.Hash
	for _, tx := range decode {
		if tx != nil {
			hash := tx.Hash()
			p.KnownTx(hash)
			if pool.AddTx(hash, tx) {
				fmt.Println("TransactionsMsg", hash, "from", p.Peer().RemoteAddr())
				hashes = append(hashes, hash)
				p.TxSync()
			}
		}
	}

	if len(hashes) > 0 {
		for _, hash := range hashes {
			backend.TxAnnounce(hash)
		}
	}
	return nil
}

func NewPooledTransactionHashesMsg66(p *peer.Peer, msg p2p.Msg) error {
	var decode packet.NewPooledTransactionHashesPacket
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("NewPooledTransactionHashesMsg: %v", err)
	}

	// size := len(decode)
	// if size >= 4096 {
	// 	return errors.New("NewPooledTransactionHashesMsg: invalid request size")
	// }

	for _, hash := range decode {
		p.KnownTx(hash)
		if pool.AddReqTx(hash) {
			p.Receiver().QueueTx(hash)
		}
	}
	return nil
}

func NewPooledTransactionHashesMsg68(p *peer.Peer, msg p2p.Msg) error {
	var decode packet.NewPooledTransactionHashesPacket68
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("NewPooledTransactionHashesMsg68: %v", err)
	}

	// if len(decode.Sizes) >= 4096 {
	// 	return errors.New("NewPooledTransactionHashesMsg68: invalid request size")
	// }
	if len(decode.Hashes) != len(decode.Types) || len(decode.Hashes) != len(decode.Sizes) {
		return fmt.Errorf("NewPooledTransactionHashesMsg68: invalid len of fields: %v %v %v", len(decode.Hashes), len(decode.Types), len(decode.Sizes))
	}

	for _, hash := range decode.Hashes {
		p.KnownTx(hash)
		if pool.AddReqTx(hash) {
			p.Receiver().QueueTx(hash)
		}
	}
	return nil
}

func NewPooledTransactionHashesMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	switch version {
	case 68:
		return NewPooledTransactionHashesMsg68(p, msg)
	default:
		return NewPooledTransactionHashesMsg66(p, msg)
	}
}

func GetPooledTransactionsMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.GetPooledTransactionsPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("GetPooledTransactionsMsg: %v", err)
	}

	var txs []*types.Transaction
	for _, hash := range decode.GetPooledTransactionsPacket {
		if tx := pool.GetTx(hash); tx != nil && tx.Done() {
			txs = append(txs, tx.Tx())
		}
	}

	if len(txs) > 0 {
		p.SendPooledTxs(txs, decode.RequestId)
	}
	return nil
}

func PooledTransactionsMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.PooledTransactionsPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("PooledTransactionsMsg: %v", err)
	}

	var hashes []common.Hash
	if p.IsRequest(decode.RequestId) {
		for _, tx := range decode.PooledTransactionsPacket {
			if tx != nil {
				hash := tx.Hash()
				if pool.UpdateTx(hash, tx) {
					fmt.Println("PooledTransactionsMsg", hash, "from", p.Peer().RemoteAddr())
					hashes = append(hashes, hash)
					p.TxSync()
				}
			}
		}
	}

	if len(hashes) > 0 {
		for _, hash := range hashes {
			backend.TxAnnounce(hash)
		}
	}
	return nil
}
