package server

import (
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/net"
	"github.com/paulhobbel/performcraft/pkg/proto"
	"github.com/paulhobbel/performcraft/pkg/pubsub"
	"log"
)

type Server struct {
	def proto.ProtocolDefinition

	publisher pubsub.Publisher
	listener  *net.Listener

	close chan bool
}

func NewServer(ver common.ProtocolVersion) *Server {
	def := proto.Registry.GetDefinition(ver)
	if def == nil {
		panic("failed to find definition, was it registered?")
	}

	return &Server{
		def:       def,
		publisher: pubsub.NewPublisher(),
		close:     make(chan bool, 1),
	}
}

func (s *Server) Listen(port int) (err error) {
	s.listener, err = net.ListenMC(fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed starting server: %w", err)
	}

	//go func() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handleConn(&conn)
	}
	//}()
}

func (s *Server) SubscribeAs(handler interface{}) {
	s.publisher.SubscribeAs(handler)
}

func (s *Server) handleConn(conn *net.Conn) {
	defer conn.Close()

	conn.Decoder.SetPacketFactory(s.def.GetClientPacket)

	for {
		p, err := conn.ReadPacket()
		if err != nil {
			break
		}

		s.publisher.PublishAs(p)
		s.publisher.PublishAs(p, conn)
	}
}
