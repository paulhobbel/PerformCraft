package proto

import (
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/util"
)

type ClientPacketStatusRequest struct {
}

func (p ClientPacketStatusRequest) ID() base.PacketID {
	return StatusRequest
}

func (p *ClientPacketStatusRequest) Read(b util.ByteBuffer) error {
	return nil
}

func (p ClientPacketStatusRequest) Write(b util.ByteBuffer) error {
	return nil
}

func (p ClientPacketStatusRequest) String() string {
	return "ClientPacketStatusRequest{}"
}
