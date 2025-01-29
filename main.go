package main

import (
	"net/http"
	_ "net/http/pprof"

	"0xlayer/go-layer-mesh/eth"
	"0xlayer/go-layer-mesh/utils/gopool"
)

func main() {
	// pprof server
	gopool.Submit(func() {
		http.ListenAndServe(":6060", nil)
	})

	server := eth.Create("bsc-mainnet")
	server.Start()
	<-server.Close()
}
