package peer

import (
	"errors"
	"fmt"
	"math/big"

	"0xlayer/go-layer-mesh/eth/message"
	"0xlayer/go-layer-mesh/eth/message/packet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
)

func (p *Peer) isStatusMsg(status *packet.StatusPacket) error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return fmt.Errorf("failed to read first message: %v", err)
	}

	if msg.Code != message.StatusMsg {
		return errors.New("invalid status message")
	}

	if msg.Size > message.MaxMessageSize {
		return errors.New("invalid status message size")
	}

	if err := msg.Decode(&status); err != nil {
		return err
	}
	return nil
}

func (p *Peer) isBroadcastTx() error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return fmt.Errorf("failed to read second message: %v", err)
	}

	if msg.Code != message.UpgradeStatusMsg {
		return errors.New("invalid upgrade status message")
	}

	if msg.Size > message.MaxMessageSize {
		return errors.New("invalid upgrade status message size")
	}

	var upgradeStatus packet.UpgradeStatusPacket
	if err := msg.Decode(&upgradeStatus); err != nil {
		return err
	}

	extension, err := upgradeStatus.GetExtension()
	if err != nil {
		return err
	}

	if extension.DisablePeerTxBroadcast {
		if !p.IsTrusted() {
			return errors.New("peer tx broadcast is disabled")
		}
		p.ext.BroadcastTx = false
	}
	return nil
}

func (p *Peer) Handshake(version uint32, networkId uint64) error {
	// todo: status from public rpc
	if p.IsTrusted() {
		if err := p.Send(message.StatusMsg, &packet.StatusPacket{
			ProtocolVersion: version,
			NetworkID:       networkId,
			TD:              big.NewInt(0),
			Head:            common.Hash{},
			Genesis:         common.Hash{},
			ForkID:          forkid.ID{},
		}); err != nil {
			return err
		}
	}

	var pStatus packet.StatusPacket
	if err := p.isStatusMsg(&pStatus); err != nil {
		return err
	}

	if pStatus.ProtocolVersion != version || pStatus.NetworkID != networkId {
		return errors.New("protocol version or network id is not matched")
	}

	if !p.IsTrusted() {
		params := &packet.StatusPacket{
			ProtocolVersion: pStatus.ProtocolVersion,
			NetworkID:       pStatus.NetworkID,
			TD:              pStatus.TD,
			Head:            pStatus.Head,
			Genesis:         pStatus.Genesis,
			ForkID:          pStatus.ForkID,
		}
		if err := p.Send(message.StatusMsg, params); err != nil {
			return err
		}
	}

	if version >= 67 {
		ext := &packet.UpgradeStatusExtension{
			DisablePeerTxBroadcast: false,
		}
		raw, err := ext.Encode()
		if err != nil {
			return err
		}

		if err := p.Send(message.UpgradeStatusMsg, &packet.UpgradeStatusPacket{
			Extension: raw,
		}); err != nil {
			return err
		}

		if err := p.isBroadcastTx(); err != nil {
			return err
		}
	}

	p.info = &PeerInfo{
		Name:    p.Peer().Info().Name,
		Enode:   p.Peer().Info().Enode,
		Version: uint(pStatus.ProtocolVersion),
	}
	fmt.Println(p.info)
	return nil
}
