package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/util"
)

type ClientPacketHandshake struct {
	Version int32
	Host    string
	Port    int16
	State   int32
}

func (p ClientPacketHandshake) ID() base.PacketID {
	return HandshakingHandshake
}

func (p *ClientPacketHandshake) Read(b util.ByteBuffer) (err error) {
	p.Version, err = b.ReadVarInt()
	p.Host, err = b.ReadString()
	p.Port, err = b.ReadShort()

	p.State, err = b.ReadVarInt()

	return
}

func (p ClientPacketHandshake) Write(b util.ByteBuffer) (err error) {
	err = b.WriteVarInt(p.Version)
	err = b.WriteString(p.Host)
	err = b.WriteShort(p.Port)
	err = b.WriteVarInt(p.State)

	return
}

func (p ClientPacketHandshake) String() string {
	return fmt.Sprintf("ClientPacketHandshake{Version: %v, Host: %v, Port: %v, NextState: %v}",
		p.Version, p.Host, p.Port, p.State)
}
