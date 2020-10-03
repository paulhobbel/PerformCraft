package r578

import (
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
)

type ClientPacketLoginStart struct {
	Name string
}

func (ClientPacketLoginStart) ID() common.PacketID {
	return LoginStart
}

func (p *ClientPacketLoginStart) Read(b common.Buffer) (err error) {
	p.Name, err = b.ReadString()

	return
}

func (p ClientPacketLoginStart) Write(b common.Buffer) error {
	return b.WriteString(p.Name)
}

func (p ClientPacketLoginStart) String() string {
	return fmt.Sprintf("ClientPacketLoginStart{Name: %s}", p.Name)
}
