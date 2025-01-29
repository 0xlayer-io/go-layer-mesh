package eth

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"0xlayer/go-layer-mesh/eth/backend"
	"0xlayer/go-layer-mesh/eth/network"
	"0xlayer/go-layer-mesh/p2p"
	"0xlayer/go-layer-mesh/p2p/enode"
	"0xlayer/go-layer-mesh/p2p/nat"
	"0xlayer/go-layer-mesh/utils"
)

type Server struct {
	s *p2p.Server
	c chan struct{}
}

func Name() string {
	// go-layer-mesh/v0.0.1/linux-amd64/go1.20.5
	return fmt.Sprintf("go-layer-mesh/v0.0.1/%s-%s/%s", runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func Create(name string) *Server {
	chain := network.GetChain(name)
	if chain == nil {
		panic(fmt.Sprintf("chain '%s' not found", name))
	}

	var (
		NAT        nat.Interface
		instanceID = os.Getenv("FLY_REGION")
		privateKey = utils.LoadOrGenerateKey(os.Getenv(fmt.Sprintf("NODE_PRIVATE_KEY_%s", instanceID)))
		maxPeers   = utils.StringToInt(os.Getenv("NODE_MAX_PEERS"), 100)
		externalIP = utils.StringOr(os.Getenv("NODE_EXTERNAL_IP"), "0.0.0.0")
		port       = utils.StringToInt(os.Getenv("NODE_PORT"), 30311)
	)

	if externalIP != "0.0.0.0" && externalIP != "127.0.0.1" {
		ip := net.ParseIP(externalIP)
		if ip == nil {
			panic(fmt.Errorf("invalid external ip: %v", externalIP))
		}
		NAT = nat.ExtIP(ip)
	} else {
		NAT = nat.Any()
	}

	var trustedNodes []*enode.Node
	if nodes := utils.StringToStrs(os.Getenv("TRUSTED_NODES")); len(nodes) > 0 {
		trustedNodes = utils.ParseNodes(nodes)
	}

	server := p2p.Server{
		Config: p2p.Config{
			Name:           Name(),
			PrivateKey:     privateKey,
			ListenAddr:     fmt.Sprintf(":%d", port),
			NAT:            NAT,
			Protocols:      network.Protocols(chain, trustedNodes),
			BootstrapNodes: chain.BootstrapNodes,
			TrustedNodes:   trustedNodes,
			// NodeDatabase:    "db/nodes", // leveldb, memorydb
			MaxPeers:        (maxPeers + len(chain.BootstrapNodes)),
			MaxPendingPeers: 10,
			DialRatio:       2,
			DiscoveryV4:     true,
			DiscoveryV5:     true,
		},
	}

	return &Server{
		s: &server,
		c: make(chan struct{}),
	}
}

func (s *Server) Start() {
	defer backend.Start()
	if err := s.s.Start(); err != nil {
		panic(err)
	}
	utils.PrintJson(s.s.NodeInfo())
	for _, n := range s.s.TrustedNodes {
		s.s.AddPeer(n)
	}
}

func (s *Server) Stop() {
	defer close(s.c)
	s.s.Stop()
}

func (s *Server) Close() <-chan struct{} {
	return s.c
}
