package utils

import "0xlayer/go-layer-mesh/p2p/enode"

func ParseNode(val string) *enode.Node {
	node, err := enode.ParseV4(val)
	if err != nil {
		panic(err)
	}
	return node
}

func ParseNodes(vals []string) []*enode.Node {
	var nodes []*enode.Node
	for _, v := range vals {
		nodes = append(nodes, ParseNode(v))
	}
	return nodes
}
