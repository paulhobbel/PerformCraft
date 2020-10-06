package network

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/pubsub"
	"io"
	"log"
	"net"
)

type Network interface {
	GetSession(id uuid.UUID) (Session, bool)
	Subscribe(handler interface{})

	Listen(port int) error

	base.Tickable
	io.Closer
}

type networkImpl struct {
	listener net.Listener

	publisher pubsub.Publisher
	sessions  map[uuid.UUID]*sessionImpl
}

func NewNetwork() Network {
	return &networkImpl{
		publisher: pubsub.NewPublisher(),
		sessions:  make(map[uuid.UUID]*sessionImpl),
	}
}

func (n networkImpl) GetSession(id uuid.UUID) (Session, bool) {
	session, found := n.sessions[id]
	return session, found
}

func (n networkImpl) Subscribe(handler interface{}) {
	n.publisher.SubscribeAs(handler)
}

func (n *networkImpl) Listen(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("cannot bind to port %d: %w", port, err)
	}
	n.listener = listener

	log.Printf("[Network]: Listening to port %d", port)

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

		session := newRemoteSession(uuid.New(), conn, n)

		n.sessions[session.id] = session

		session.init()
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
