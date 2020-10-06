package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/util"
)

type ClientPacketLoginStart struct {
	Name string
}

func (ClientPacketLoginStart) ID() base.PacketID {
	return LoginStart
}

func (p *ClientPacketLoginStart) Read(b util.ByteBuffer) (err error) {
	p.Name, err = b.ReadString()

	return
}

func (p ClientPacketLoginStart) Write(b util.ByteBuffer) error {
	return b.WriteString(p.Name)
}

func (p ClientPacketLoginStart) String() string {
	return fmt.Sprintf("ClientPacketLoginStart{Name: %s}", p.Name)
}
