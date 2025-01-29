package packet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type GetReceiptsPacket []common.Hash
type GetReceiptsPacket66 struct {
	RequestId uint64
	GetReceiptsPacket
}

type ReceiptsPacket [][]*types.Receipt
type ReceiptsPacket66 struct {
	RequestId uint64
	ReceiptsPacket
}

type ReceiptsRLPPacket []rlp.RawValue
type ReceiptsRLPPacket66 struct {
	RequestId uint64
	ReceiptsRLPPacket
}
