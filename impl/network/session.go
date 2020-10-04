package network

import (
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/game"
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
	"net"
)

type sessionImpl struct {
	id      uuid.UUID
	profile *game.Profile

	threshold int

	protocolVersion proto.ProtocolVersion
	protocolState   proto.ProtocolState

	conn    net.Conn
	network *networkImpl
}

func (s sessionImpl) GetId() uuid.UUID {
	return s.id
}

func (s sessionImpl) GetPlayerProfile() *game.Profile {
	return s.profile
}

func (s *sessionImpl) SetThreshold(threshold int) {
	s.threshold = threshold
}

func (s *sessionImpl) SetProtocolVersion(version proto.ProtocolVersion) {
	s.protocolVersion = version
}

func (s *sessionImpl) SetProtocolState(state proto.ProtocolState) {
	s.protocolState = state
}

func (sessionImpl) WritePacket(packet network.Packet) {
	panic("implement me")
}

func (sessionImpl) Tick(deltaTicks int) {
	panic("implement me")
}

func (s sessionImpl) readPipe() {
	defer s.Close()
	for {

	}
}

func (s sessionImpl) Close() error {
	delete(s.network.sessions, s.id)

	return s.conn.Close()
}
