package proto

import (
	"github.com/paulhobbel/performcraft/pkg/common"
	"log"
)

type ProtocolVersion int

const (
	R578 ProtocolVersion = 578
	R753 ProtocolVersion = 753
)

type packetMapping map[common.PacketState]map[common.PacketID]func() common.Packet

type ProtocolDefinition interface {
	GetClientPacket(state common.PacketState, id common.PacketID) common.Packet
	GetServerPacket(state common.PacketState, id common.PacketID) common.Packet
}

type protocolDefinitionImpl struct {
	clientPackets packetMapping
	serverPackets packetMapping
}

func (def protocolDefinitionImpl) GetClientPacket(state common.PacketState, id common.PacketID) common.Packet {
	return def.getSidedPacket(def.clientPackets, state, id)
}

func (def protocolDefinitionImpl) GetServerPacket(state common.PacketState, id common.PacketID) common.Packet {
	return def.getSidedPacket(def.serverPackets, state, id)
}

func (def protocolDefinitionImpl) getSidedPacket(mapping packetMapping, state common.PacketState, id common.PacketID) common.Packet {
	factory := mapping[state][id]
	if factory == nil {
		return nil
	}

	return factory()
}

var Registry = &protocolRegistry{map[ProtocolVersion]ProtocolDefinition{}}

type protocolRegistry struct {
	definitions map[ProtocolVersion]ProtocolDefinition
}

func (reg *protocolRegistry) RegisterProtocol(ver ProtocolVersion) *ProtocolDefinitionBuilder {
	def := &protocolDefinitionImpl{
		clientPackets: make(packetMapping),
		serverPackets: make(packetMapping),
	}

	reg.definitions[ver] = def
	log.Printf("[ProtocolRegistry] Registered %d protocol", ver)

	return &ProtocolDefinitionBuilder{def}
}

func (reg *protocolRegistry) GetDefinition(ver ProtocolVersion) ProtocolDefinition {
	def := reg.definitions[ver]
	return def
}
