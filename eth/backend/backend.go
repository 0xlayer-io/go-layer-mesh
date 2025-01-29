package backend

import "0xlayer/go-layer-mesh/eth/backend/pool"

func Start() {
	AnnounceLoop()
	pool.BlockLoop()
	pool.TxLoop()
}
