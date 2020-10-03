package r578

import "github.com/paulhobbel/performcraft/pkg/common"

type ClientPacketStatusRequest struct {
}

func (p ClientPacketStatusRequest) ID() common.PacketID {
	return StatusRequest
}

func (p *ClientPacketStatusRequest) Read(b common.Buffer) error {
	return nil
}

func (p ClientPacketStatusRequest) Write(b common.Buffer) error {
	return nil
}

func (p ClientPacketStatusRequest) String() string {
	return "ClientPacketStatusRequest{}"
}
