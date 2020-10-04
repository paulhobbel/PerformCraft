package handler

import (
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/net"
	"github.com/paulhobbel/performcraft/pkg/proto/r578"
	"log"
)

func HandshakeHandler(handshake *r578.ClientPacketHandshake, conn *net.Conn) {
	nextState := common.ProtocolStateFromId(handshake.State)
	if nextState == common.Unknown {
		log.Printf("[HandshakeHandler]: Unkown protocol state! (%d)", handshake.State)
		// TODO: Add disconnect function
		conn.Close()
		return
	}

	conn.SetProtocolState(nextState)
	if nextState != common.Status && nextState != common.Login {
		log.Printf("[HandshakeHandler]: Unexpected next protocol state! (%d), expected Status or Login", handshake.State)
		// TODO: Add disconnect function
		conn.Close()
		return
	}

}
