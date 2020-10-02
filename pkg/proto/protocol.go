package proto

import (
	"github.com/paulhobbel/performcraft/pkg/common"
)

type ProtocolVersion int

const (
	R578 ProtocolVersion = 578
)

type ProtocolDefinition interface {
	GetClientPacket(state common.PacketState, id common.PacketID) interface{}
}
