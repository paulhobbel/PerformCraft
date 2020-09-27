package net

import (
	"bufio"
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
		Socket: conn,
		Reader: bufio.NewReader(conn),
		Writer: conn,
	}, err
}
