package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/util"
)

type PacketStatusPing struct {
	Payload int64
}

func (p PacketStatusPing) ID() base.PacketID {
	return StatusPing
}

func (p *PacketStatusPing) Read(b util.ByteBuffer) (err error) {
	p.Payload, err = b.ReadLong()

	return
}

func (p PacketStatusPing) Write(b util.ByteBuffer) error {
	return b.WriteLong(p.Payload)
}

func (p PacketStatusPing) String() string {
	return fmt.Sprintf("PacketStatusPing{Payload: %v}", p.Payload)
}
