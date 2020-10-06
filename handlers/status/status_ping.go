package status

import (
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
)

func StatusPingHandler(request *proto.PacketStatusPing, session network.Session) {
	session.WritePacket(request)
}
