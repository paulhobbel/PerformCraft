package proto

import (
	"github.com/paulhobbel/performcraft/pkg/common"
)

type ProtocolDefinitionBuilder struct {
	def *protocolDefinitionImpl
}

func (b *ProtocolDefinitionBuilder) RegisterClientPacket(state common.PacketState, id common.PacketID, packet common.Packet) *ProtocolDefinitionBuilder {
	if b.def.clientPackets[state] == nil {
		b.def.clientPackets[state] = make(map[common.PacketID]func() common.Packet)
	}

	b.def.clientPackets[state][id] = func() common.Packet {
		return packet
	}

	return b
}

func (b *ProtocolDefinitionBuilder) RegisterServerPacket(state common.PacketState, id common.PacketID, packet common.Packet) *ProtocolDefinitionBuilder {
	if b.def.serverPackets[state] == nil {
		b.def.serverPackets[state] = make(map[common.PacketID]func() common.Packet)
	}

	b.def.serverPackets[state][id] = func() common.Packet {
		return packet
	}

	return b
}
