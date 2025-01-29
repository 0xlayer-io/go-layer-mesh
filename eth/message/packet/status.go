package packet

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/rlp"
)

type StatusPacket struct {
	ProtocolVersion uint32
	NetworkID       uint64
	TD              *big.Int
	Head            common.Hash
	Genesis         common.Hash
	ForkID          forkid.ID
}

type UpgradeStatusExtension struct {
	DisablePeerTxBroadcast bool
}

func (e *UpgradeStatusExtension) Encode() (*rlp.RawValue, error) {
	rawBytes, err := rlp.EncodeToBytes(e)
	if err != nil {
		return nil, err
	}
	raw := rlp.RawValue(rawBytes)
	return &raw, nil
}

type UpgradeStatusPacket struct {
	Extension *rlp.RawValue `rlp:"nil"`
}

func (p *UpgradeStatusPacket) GetExtension() (*UpgradeStatusExtension, error) {
	extension := &UpgradeStatusExtension{}
	if p.Extension == nil {
		return extension, nil
	}
	err := rlp.DecodeBytes(*p.Extension, extension)
	if err != nil {
		return nil, err
	}
	return extension, nil
}
