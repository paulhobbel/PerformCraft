package net

import (
	"bufio"
	"github.com/paulhobbel/performcraft/pkg/net/packet"
	"io"
	"net"
)

type Conn struct {
	Socket net.Conn
	Reader packet.Reader

	io.Writer

	threshold int
}

func DialMC(addr string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	return &Conn{
		Socket: conn,
		Reader: bufio.NewReader(conn),
		Writer: conn,
	}, err
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	p, err := packet.ReadPacket(c.Reader, c.threshold > 0)
	if err != nil {
		return packet.Packet{}, err
	}
	return *p, err
}

func (c *Conn) WritePacket(p packet.Packet) error {
	_, err := c.Write(p.Pack(c.threshold))
	return err
}

func (c *Conn) Close() error {
	return c.Socket.Close()
}

func (c *Conn) SetThreshold(threshold int) {
	c.threshold = threshold
}
