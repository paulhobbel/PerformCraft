package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
)

type BasePacket struct {
	id   base.PacketID
	data []byte
}

func (p BasePacket) ID() base.PacketID {
	return p.id
}

func (p *BasePacket) Read(b bufio.ByteBuffer) error {
	p.data = b.Bytes()
	return nil
}

func (p BasePacket) Write(b bufio.ByteBuffer) error {
	_, err := b.Write(p.data)

	return err
}

func (p BasePacket) String() string {
	return fmt.Sprintf("Packet{ID: 0x%02x, Data: %v}", p.id, p.data)
}
