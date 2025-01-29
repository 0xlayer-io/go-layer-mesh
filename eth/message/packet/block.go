package packet

import (
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type HashOrNumber struct {
	Hash   common.Hash
	Number uint64
}

type BlobSidecars []*BlobSidecar
type BlobSidecar struct {
	BlobTxSidecar *types.BlobTxSidecar
	BlockNumber   *big.Int    `json:"blockNumber"`
	BlockHash     common.Hash `json:"blockHash"`
	TxIndex       uint64      `json:"transactionIndex"`
	TxHash        common.Hash `json:"transactionHash"`
}

type BlockBody struct {
	Transactions []*types.Transaction
	Uncles       []*types.Header
	Withdrawals  []*types.Withdrawal `rlp:"optional"`
	Sidecars     BlobSidecars        `rlp:"optional"`
}

type NewBlockPacket struct {
	Block    *types.Block
	TD       *big.Int
	Sidecars BlobSidecars `rlp:"optional"`
}

type NewBlockHashesPacket []struct {
	Hash   common.Hash
	Number uint64
}

type GetBlockHeadersPacket struct {
	Origin  HashOrNumber
	Amount  uint64
	Skip    uint64
	Reverse bool
}
type GetBlockHeadersPacket66 struct {
	RequestId uint64
	*GetBlockHeadersPacket
}

type BlockHeadersPacket []*types.Header
type BlockHeadersPacket66 struct {
	RequestId uint64
	BlockHeadersPacket
}

type BlockHeadersRLPPacket []rlp.RawValue
type BlockHeadersRLPPacket66 struct {
	RequestId uint64
	BlockHeadersRLPPacket
}

type GetBlockBodiesPacket []common.Hash
type GetBlockBodiesPacket66 struct {
	RequestId uint64
	GetBlockBodiesPacket
}

type BlockBodiesPacket []*BlockBody
type BlockBodiesPacket66 struct {
	RequestId uint64
	BlockBodiesPacket
}

type BlockBodiesRLPPacket []rlp.RawValue
type BlockBodiesRLPPacket66 struct {
	RequestId uint64
	BlockBodiesRLPPacket
}

func (p *NewBlockHashesPacket) Unpack() ([]common.Hash, []uint64) {
	var (
		hashes  = make([]common.Hash, len(*p))
		numbers = make([]uint64, len(*p))
	)
	for i, body := range *p {
		hashes[i], numbers[i] = body.Hash, body.Number
	}
	return hashes, numbers
}

func (p *BlockBodiesPacket) Unpack() ([][]*types.Transaction, [][]*types.Header, [][]*types.Withdrawal, []BlobSidecars) {
	var (
		txset         = make([][]*types.Transaction, len(*p))
		uncleset      = make([][]*types.Header, len(*p))
		withdrawalset = make([][]*types.Withdrawal, len(*p))
		sidecarset    = make([]BlobSidecars, len(*p))
	)
	for i, body := range *p {
		txset[i], uncleset[i], withdrawalset[i], sidecarset[i] = body.Transactions, body.Uncles, body.Withdrawals, body.Sidecars
	}
	return txset, uncleset, withdrawalset, sidecarset
}

func (hn *HashOrNumber) EncodeRLP(w io.Writer) error {
	if hn.Hash == (common.Hash{}) {
		return rlp.Encode(w, hn.Number)
	}
	if hn.Number != 0 {
		return fmt.Errorf("both origin hash (%x) and number (%d) provided", hn.Hash, hn.Number)
	}
	return rlp.Encode(w, hn.Hash)
}

func (hn *HashOrNumber) DecodeRLP(s *rlp.Stream) error {
	_, size, err := s.Kind()
	switch {
	case err != nil:
		return err
	case size == 32:
		hn.Number = 0
		return s.Decode(&hn.Hash)
	case size <= 8:
		hn.Hash = common.Hash{}
		return s.Decode(&hn.Number)
	default:
		return fmt.Errorf("invalid input size %d for origin", size)
	}
}
