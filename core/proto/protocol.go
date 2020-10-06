package proto

import (
	"github.com/paulhobbel/performcraft/core/base"
)

var clientPackets map[base.ProtocolState]map[base.PacketID]func() base.Packet

func GetVersion() base.ProtocolVersion {
	return base.ProtocolVersion(578)
}

func GetClientPacket(state base.ProtocolState, packetId base.PacketID) base.Packet {
	factory := clientPackets[state][packetId]
	if factory == nil {
		return &BasePacket{packetId, nil}
	}

	return factory()
}
