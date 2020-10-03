package v2

import (
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
)

type BasePacket struct {
	id   common.PacketID
	data []byte
}

func (p BasePacket) ID() common.PacketID {
	return p.id
}

func (p *BasePacket) Read(b common.Buffer) error {
	p.data = b.Bytes()
	return nil
}

func (p BasePacket) Write(b common.Buffer) error {
	_, err := b.Write(p.data)

	return err
}

func (p BasePacket) String() string {
	return fmt.Sprintf("Packet{ID: 0x%02x, Data: %v}", p.id, p.data)
}
