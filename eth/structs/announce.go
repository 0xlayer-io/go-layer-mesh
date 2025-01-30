package structs

import "github.com/ethereum/go-ethereum/common"

type BlockAnnounce struct {
	Hash   common.Hash
	Number uint64
}
