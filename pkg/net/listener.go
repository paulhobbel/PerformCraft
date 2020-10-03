package net

import (
	v2 "github.com/paulhobbel/performcraft/pkg/net/packet/v2"
	"net"
)

type Listener struct {
	net.Listener
}

func ListenMC(addr string) (*Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Listener{l}, nil
}

func (l Listener) Accept() (Conn, error) {
	conn, err := l.Listener.Accept()
	return Conn{
		Socket:  conn,
		Encoder: v2.NewEncoder(conn),
		Decoder: v2.NewDecoder(conn),
	}, err
}
