package proto

import (
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
)

type ClientPacketStatusRequest struct {
}

func (p ClientPacketStatusRequest) ID() base.PacketID {
	return StatusRequest
}

func (p *ClientPacketStatusRequest) Read(b bufio.ByteBuffer) error {
	return nil
}

func (p ClientPacketStatusRequest) Write(b bufio.ByteBuffer) error {
	return nil
}

func (p ClientPacketStatusRequest) String() string {
	return "ClientPacketStatusRequest{}"
}
