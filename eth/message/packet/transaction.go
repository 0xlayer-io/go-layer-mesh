package packet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type TransactionsPacket []*types.Transaction

type NewPooledTransactionHashesPacket []common.Hash
type NewPooledTransactionHashesPacket68 struct {
	Types  []byte
	Sizes  []uint32
	Hashes []common.Hash
}

type GetPooledTransactionsPacket []common.Hash
type GetPooledTransactionsPacket66 struct {
	RequestId uint64
	GetPooledTransactionsPacket
}

type PooledTransactionsPacket []*types.Transaction
type PooledTransactionsPacket66 struct {
	RequestId uint64
	PooledTransactionsPacket
}

type PooledTransactionsRLPPacket []rlp.RawValue
type PooledTransactionsRLPPacket66 struct {
	RequestId uint64
	PooledTransactionsRLPPacket
}
