package handler

import (
	"fmt"

	"0xlayer/go-layer-mesh/eth/handler/peer"
	"0xlayer/go-layer-mesh/eth/message"
)

func peerMsg(p *peer.Peer, version uint32) error {
	msg, err := p.Msg().ReadMsg()
	if err != nil {
		return err
	}
	defer msg.Discard()

	if msg.Size > message.MaxMessageSize {
		return fmt.Errorf("message too large: %d", msg.Size)
	}

	if ok, err := Message(p, msg, version); ok {
		if err != nil {
			fmt.Println(err, "code", msg.Code, "size", msg.Size, "from", p.Peer().RemoteAddr())
			// todo: trusted, best (tx)
			return err
		}
	} else {
		fmt.Println("Unknown message from", p.Peer().RemoteAddr(), "code", msg.Code, "size", msg.Size)
	}
	return nil
}

func Peer(p *peer.Peer, version uint32) error {
	for {
		err := peerMsg(p, version)
		if err != nil {
			return err
		}
	}
}
