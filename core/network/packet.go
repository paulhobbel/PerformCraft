package network

import (
	"fmt"
)

type PacketID int32

type Packet interface {
	ID() PacketID
	Read(buffer Buffer) error
	Write(buffer Buffer) error
}

type BasePacket struct {
	id   PacketID
	data []byte
}

func (p BasePacket) ID() PacketID {
	return p.id
}

func (p *BasePacket) Read(b Buffer) error {
	p.data = b.Bytes()
	return nil
}

func (p BasePacket) Write(b Buffer) error {
	_, err := b.Write(p.data)

	return err
}

func (p BasePacket) String() string {
	return fmt.Sprintf("Packet{ID: 0x%02x, Data: %v}", p.id, p.data)
}
