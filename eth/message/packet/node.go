package packet

import "github.com/ethereum/go-ethereum/common"

type GetNodeDataPacket []common.Hash
type GetNodeDataPacket66 struct {
	RequestId uint64
	GetNodeDataPacket
}

type NodeDataPacket [][]byte
type NodeDataPacket66 struct {
	RequestId uint64
	NodeDataPacket
}
