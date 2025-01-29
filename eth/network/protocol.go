package network

import (
	"fmt"
	"time"

	"0xlayer/go-layer-mesh/eth/handler"
	"0xlayer/go-layer-mesh/eth/handler/peer"
	"0xlayer/go-layer-mesh/p2p"
	"0xlayer/go-layer-mesh/p2p/enode"
	"0xlayer/go-layer-mesh/utils"
)

type Protocol struct {
	Name    string
	Version uint32
	Length  uint64
}

var protocols = map[string][]Protocol{
	"bsc-mainnet": {
		{"eth", 66, 17},
		{"eth", 67, 17},
		{"eth", 68, 17},
	},
	"bsc-testnet": {
		{"eth", 66, 17},
		{"eth", 67, 17},
		{"eth", 68, 17},
	},
}

func isTrusted(trustedNodes []*enode.Node, p *p2p.Peer) bool {
	for _, v := range trustedNodes {
		if v.ID() == p.Node().ID() {
			return true
		}
	}
	return false
}

func Protocols(chain *Chain, trustedNodes []*enode.Node) []p2p.Protocol {
	protocol := protocols[chain.Name]
	if len(protocol) > 0 {
		var protocols []p2p.Protocol
		for _, v := range protocol {
			protocols = append(protocols, p2p.Protocol{
				Name:    v.Name,
				Version: uint(v.Version),
				Length:  v.Length,
				Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
					peer := peer.New(p, rw, isTrusted(trustedNodes, p))
					if err := utils.Timeout(func() error {
						return peer.Handshake(v.Version, chain.NetworkID)
					}, (10 * time.Second)); err != nil {
						return fmt.Errorf("handshake failed: %v", err)
					}

					peer.Open()
					defer peer.Close()
					err := handler.Peer(peer, v.Version)
					return err
				},
				PeerInfo: func(id enode.ID) interface{} {
					if p := peer.Get(id); p != nil {
						return p.Peer().Info()
					}
					return nil
				},
				NodeInfo: func() interface{} {
					return nil
				},
			})
		}
		return protocols
	}
	return nil
}
