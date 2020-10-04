package network

import (
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/common"
	"github.com/paulhobbel/performcraft/core/game"
	"github.com/paulhobbel/performcraft/core/proto"
	"io"
)

type Session interface {
	GetId() uuid.UUID

	GetPlayerProfile() *game.Profile

	SetThreshold(threshold int)
	SetProtocolVersion(version proto.ProtocolVersion)
	SetProtocolState(state proto.ProtocolState)

	WritePacket(packet Packet)

	common.Tickable
	io.Closer
}

type RemoteSession interface {
	Session

	GetNetwork() Network

	Disconnect(reason string) error
}
