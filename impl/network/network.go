package network

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
	"log"
	"net"
)

type networkImpl struct {
	listener net.Listener

	sessions map[uuid.UUID]*sessionImpl
}

func NewNetwork() network.Network {
	return &networkImpl{
		sessions: make(map[uuid.UUID]*sessionImpl),
	}
}

func (n networkImpl) GetSession(id uuid.UUID) (network.Session, bool) {
	session, found := n.sessions[id]
	return session, found
}

func (n *networkImpl) Listen(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("cannot bind to port %d: %w", port, err)
	}
	n.listener = listener

	n.acceptConnections()

	return nil
}

func (n *networkImpl) acceptConnections() {
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			log.Printf("[Network]: Failed accepting connection: %v", err)
			continue
		}

		session := &sessionImpl{
			id:            uuid.New(),
			protocolState: proto.Handshake,
			conn:          conn,
			network:       n,
		}

		n.sessions[session.id] = session

		go session.readPipe()

		log.Printf("[Network]: Accepted new client connection with id: %s", session.id)
	}
}

func (n networkImpl) Tick(deltaTicks int) {
	for _, session := range n.sessions {
		session.Tick(deltaTicks)
	}
}

func (networkImpl) Close() error {
	panic("implement me")
}
