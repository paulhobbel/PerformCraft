package proto

import (
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/util"
)

type ClientPacketKeepAlive struct {
	KeepAliveID int64
}

func (ClientPacketKeepAlive) ID() base.PacketID {
	//return PlayKeepAliveClient
	return base.PacketID(0x0F)
}

func (p *ClientPacketKeepAlive) Read(buffer util.ByteBuffer) (err error) {
	p.KeepAliveID, err = buffer.ReadLong()

	return
}

func (p ClientPacketKeepAlive) Write(buffer util.ByteBuffer) error {
	return buffer.WriteLong(p.KeepAliveID)
}
