package proto

import (
	"github.com/paulhobbel/performcraft/pkg/common"
	"log"
)

const (
	R578 common.ProtocolVersion = 578
	R753 common.ProtocolVersion = 753
)

type packetMapping map[common.ProtocolState]map[common.PacketID]func() common.Packet

type ProtocolDefinition interface {
	GetClientPacket(state common.ProtocolState, id common.PacketID) common.Packet
	GetServerPacket(state common.ProtocolState, id common.PacketID) common.Packet
}

type protocolDefinitionImpl struct {
	clientPackets packetMapping
	serverPackets packetMapping
}

func (def protocolDefinitionImpl) GetClientPacket(state common.ProtocolState, id common.PacketID) common.Packet {
	return def.getSidedPacket(def.clientPackets, state, id)
}

func (def protocolDefinitionImpl) GetServerPacket(state common.ProtocolState, id common.PacketID) common.Packet {
	return def.getSidedPacket(def.serverPackets, state, id)
}

func (def protocolDefinitionImpl) getSidedPacket(mapping packetMapping, state common.ProtocolState, id common.PacketID) common.Packet {
	factory := mapping[state][id]
	if factory == nil {
		return nil
	}

	return factory()
}

var Registry = &protocolRegistry{map[common.ProtocolVersion]ProtocolDefinition{}}

type protocolRegistry struct {
	definitions map[common.ProtocolVersion]ProtocolDefinition
}

func (reg *protocolRegistry) RegisterProtocol(ver common.ProtocolVersion) *ProtocolDefinitionBuilder {
	def := &protocolDefinitionImpl{
		clientPackets: make(packetMapping),
		serverPackets: make(packetMapping),
	}

	reg.definitions[ver] = def
	log.Printf("[ProtocolRegistry] Registered %d protocol", ver)

	return &ProtocolDefinitionBuilder{def}
}

func (reg *protocolRegistry) GetDefinition(ver common.ProtocolVersion) ProtocolDefinition {
	def := reg.definitions[ver]
	return def
}
