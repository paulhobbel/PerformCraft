package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/text"
	"github.com/paulhobbel/performcraft/core/util"
)

type ServerPacketLoginDisconnect struct {
	Reason text.Message
}

func (p ServerPacketLoginDisconnect) ID() base.PacketID {
	return LoginDisconnect
}

func (p *ServerPacketLoginDisconnect) Read(b util.ByteBuffer) (err error) {
	p.Reason, err = b.ReadMessage()

	return
}

func (p ServerPacketLoginDisconnect) Write(b util.ByteBuffer) error {
	return b.WriteMessage(p.Reason)
}

func (p ServerPacketLoginDisconnect) String() string {
	return fmt.Sprintf("ServerPacketLoginDisconnect{Reason: %v}", p.Reason)
}

type ServerPacketLoginSuccess struct {
	UUID string
	Name string
}

func (ServerPacketLoginSuccess) ID() base.PacketID {
	return LoginSuccess
}

func (p *ServerPacketLoginSuccess) Read(b util.ByteBuffer) (err error) {
	p.UUID, err = b.ReadString()
	p.Name, err = b.ReadString()

	return
}

func (p ServerPacketLoginSuccess) Write(b util.ByteBuffer) (err error) {
	err = b.WriteString(p.UUID)
	err = b.WriteString(p.Name)

	return
}

func (p ServerPacketLoginSuccess) String() string {
	return fmt.Sprintf("ServerPacketLoginSuccess{UUID: %s, Name: %s", p.UUID, p.Name)
}

type ServerPacketSetCompression struct {
	Threshold int32
}

func (p ServerPacketSetCompression) ID() base.PacketID {
	return LoginSetCompression
}

func (p *ServerPacketSetCompression) Read(b util.ByteBuffer) (err error) {
	p.Threshold, err = b.ReadVarInt()

	return
}

func (p ServerPacketSetCompression) Write(b util.ByteBuffer) error {
	return b.WriteVarInt(p.Threshold)
}

func (p ServerPacketSetCompression) String() string {
	return fmt.Sprintf("ServerPacketSetCompression{Threshold: %d}", p.Threshold)
}
