package r578

import (
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
)

type PacketStatusPing struct {
	Payload int64
}

func (p PacketStatusPing) ID() common.PacketID {
	return StatusPing
}

func (p *PacketStatusPing) Read(b common.Buffer) (err error) {
	p.Payload, err = b.ReadLong()

	return
}

func (p PacketStatusPing) Write(b common.Buffer) error {
	return b.WriteLong(p.Payload)
}

func (p PacketStatusPing) String() string {
	return fmt.Sprintf("PacketStatusPing{Payload: %v}", p.Payload)
}
