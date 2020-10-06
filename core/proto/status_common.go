package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
)

type PacketStatusPing struct {
	Payload int64
}

func (p PacketStatusPing) ID() base.PacketID {
	return StatusPing
}

func (p *PacketStatusPing) Read(b bufio.ByteBuffer) (err error) {
	p.Payload, err = b.ReadLong()

	return
}

func (p PacketStatusPing) Write(b bufio.ByteBuffer) error {
	return b.WriteLong(p.Payload)
}

func (p PacketStatusPing) String() string {
	return fmt.Sprintf("PacketStatusPing{Payload: %v}", p.Payload)
}
