package handshake

import (
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
	"log"
)

func HandshakeHandler(request *proto.ClientPacketHandshake, session network.Session) {
	nextState := base.ProtocolStateFromId(request.State)
	if nextState == base.Unknown {
		log.Printf("[Session]: Unkown protocol state! (%d)", request.State)
		// TODO: Add disconnect function
		session.Close()
		return
	}

	session.SetProtocolState(nextState)
	if nextState != base.Status && nextState != base.Login {
		log.Printf("[Session]: Unexpected next protocol state! (%d), expected Status or Login", request.State)
		// TODO: Add disconnect function
		session.Close()
		return
	}
}