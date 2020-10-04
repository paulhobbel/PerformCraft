package net

import (
	"github.com/paulhobbel/performcraft/pkg/common"
	v2 "github.com/paulhobbel/performcraft/pkg/net/packet/v2"
	"net"
)

type Conn struct {
	Socket  net.Conn
	Encoder *v2.Encoder
	Decoder *v2.Decoder

	State common.ProtocolState

	protocolVersion common.ProtocolVersion
}

func DialMC(addr string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	return &Conn{
		Socket:  conn,
		Encoder: v2.NewEncoder(conn),
		Decoder: v2.NewDecoder(conn),
	}, err
}

func (c *Conn) ReadPacket() (common.Packet, error) {
	return c.Decoder.Unmarshal(c.State)
}

func (c *Conn) WritePacket(p common.Packet) error {
	//log.Printf("[Conn]: Writing packet: %+v", p)
	return c.Encoder.Marshal(p)
}

func (c *Conn) Close() error {
	return c.Socket.Close()
}

func (c *Conn) SetThreshold(threshold int) {
	c.Encoder.SetThreshold(threshold)
	c.Decoder.SetThreshold(threshold)
}

func (c *Conn) SetProtocolVersion(version common.ProtocolVersion) {
	c.protocolVersion = version
}

func (c *Conn) SetProtocolState(state common.ProtocolState) {
	c.State = state
}
