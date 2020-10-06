package network

import (
	"errors"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/entity"
	"github.com/paulhobbel/performcraft/core/game"
	"github.com/paulhobbel/performcraft/core/proto"
	"github.com/paulhobbel/performcraft/core/pubsub"
	"io"
	"log"
	"net"
	"time"
)

type Session interface {
	GetId() uuid.UUID

	GetPlayer() entity.Player
	SetPlayerProfile(profile *game.Profile)
	GetPlayerProfile() *game.Profile

	SetThreshold(threshold int)
	SetProtocolState(state base.ProtocolState)

	WritePacket(packet base.Packet) error

	InitPlayer()

	base.Tickable
	io.Closer
}

type RemoteSession interface {
	Session

	GetNetwork() Network

	Disconnect(reason string) error
}

type sessionImpl struct {
	id      uuid.UUID
	player  entity.Player
	profile *game.Profile

	threshold int

	protocolState base.ProtocolState

	conn    net.Conn
	network *networkImpl

	encoder *packetEncoder
	decoder *packetDecoder

	publisher pubsub.Publisher

	close chan bool
}

func newRemoteSession(id uuid.UUID, conn net.Conn, network *networkImpl) *sessionImpl {
	return &sessionImpl{
		id:            id,
		profile:       nil,
		threshold:     0,
		protocolState: base.Handshake,
		conn:          conn,
		network:       network,
		encoder:       NewPacketEncoder(conn),
		decoder:       NewPacketDecoder(conn),
		publisher:     pubsub.WrapPublisher(network.publisher),
		close:         make(chan bool, 1),
	}
}

func (s sessionImpl) GetId() uuid.UUID {
	return s.id
}

func (s sessionImpl) GetPlayer() entity.Player {
	return s.player
}

func (s sessionImpl) GetPlayerProfile() *game.Profile {
	return s.profile
}

func (s *sessionImpl) SetPlayerProfile(profile *game.Profile) {
	s.profile = profile
}

func (s *sessionImpl) SetThreshold(threshold int) {
	s.threshold = threshold
	s.decoder.SetThreshold(threshold)
	s.encoder.SetThreshold(threshold)
	log.Printf("[Session]: Changed session compression threshold to %d", threshold)
}

func (s *sessionImpl) SetProtocolState(state base.ProtocolState) {
	s.protocolState = state
	s.decoder.SetProtocolState(state)
}

func (s sessionImpl) WritePacket(packet base.Packet) error {
	return s.encoder.Encode(packet)
}

func (s sessionImpl) InitPlayer() {
	s.SetProtocolState(base.Play)
	go s.initKeepAlive()

	// TODO: Create player
	s.WritePacket(&proto.ServerPacketPlayJoinGame{
		EntityId:          0,
		GameMode:          1,
		Dimension:         0,
		Seed:              0,
		MaxPlayers:        0,
		LevelType:         "default",
		ViewDistance:      15,
		ReducedDebugInfo:  false,
		EnableSpawnScreen: false,
	})

	s.WritePacket(&proto.ServerPacketPlayPlayerPositionLook{
		X:          0,
		Y:          0,
		Z:          0,
		Yaw:        0,
		Pitch:      0,
		Flags:      0,
		TeleportID: 0,
	})
}

func (sessionImpl) Tick(deltaTicks int) {
	//log.Printf("[Session]: Handling %d ticks", deltaTicks)
}

func (s sessionImpl) init() {
}

func (s *sessionImpl) initKeepAlive() {
	ticker := time.NewTicker(time.Second * 15)

	for {
		select {
		case <-s.close:
			log.Println("[Session]: Stopped keep alive loop")
			return
		case curr := <-ticker.C:
			if s.protocolState == base.Play {
				// TODO: Handle client not responding...
				log.Printf("[Session]: Sending keep alive")
				s.WritePacket(&proto.ServerPacketKeepAlive{KeepAliveID: curr.UnixNano() / 1e6})
			}
		}
	}
}

func (s *sessionImpl) readPipe() {
	defer s.Close()
	for {
		p, err := s.decoder.Decode()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("[Session]: Closing connection: %v", err)
				break
			}
		}

		if p != nil {
			log.Printf("[Session]: Received packet: %+v", p)
			s.publisher.PublishAs(p)
			s.publisher.PublishAs(p, s)
		}

	}
}

func (s *sessionImpl) Close() error {
	s.close <- true

	delete(s.network.sessions, s.id)
	close(s.close)

	return s.conn.Close()
}

func createStatusResponse() *proto.ServerPacketStatusResponse {
	info := &proto.ServerPacketStatusResponse{}

	info.Version.Name = "PerformCraft 1.15.2"
	info.Version.Protocol = 578
	info.Players.Max = 1337
	info.Players.Online = 420
	//info.Players.Sample =
	info.Description = "PerformCraft Research Server"

	return info
}
