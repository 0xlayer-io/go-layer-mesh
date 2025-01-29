package handler_message

import (
	"fmt"

	"0xlayer/go-layer-mesh/eth/handler/peer"
	"0xlayer/go-layer-mesh/eth/message/packet"
	"0xlayer/go-layer-mesh/p2p"
)

func GetReceiptsMsg(p *peer.Peer, msg p2p.Msg, version uint32) error {
	var decode packet.GetReceiptsPacket66
	if err := msg.Decode(&decode); err != nil {
		return fmt.Errorf("GetReceiptsMsg: %v", err)
	}
	p.SendEmptyNodeData(decode.RequestId)
	return nil
}
