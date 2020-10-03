package r578

import (
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/net/chat"
)

type ServerPacketLoginDisconnect struct {
	Reason chat.Message
}

func (p ServerPacketLoginDisconnect) ID() common.PacketID {
	return LoginDisconnect
}

func (p *ServerPacketLoginDisconnect) Read(b common.Buffer) (err error) {
	p.Reason, err = b.ReadMessage()

	return
}

func (p ServerPacketLoginDisconnect) Write(b common.Buffer) error {
	return b.WriteMessage(p.Reason)
}

func (p ServerPacketLoginDisconnect) String() string {
	return fmt.Sprintf("ServerPacketLoginDisconnect{Reason: %v}", p.Reason)
}

type ServerPacketLoginSuccess struct {
	UUID string
	Name string
}

func (ServerPacketLoginSuccess) ID() common.PacketID {
	return LoginSuccess
}

func (p *ServerPacketLoginSuccess) Read(b common.Buffer) (err error) {
	p.UUID, err = b.ReadString()
	p.Name, err = b.ReadString()

	return
}

func (p ServerPacketLoginSuccess) Write(b common.Buffer) (err error) {
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

func (p ServerPacketSetCompression) ID() common.PacketID {
	return LoginSetCompression
}

func (p *ServerPacketSetCompression) Read(b common.Buffer) (err error) {
	p.Threshold, err = b.ReadVarInt()

	return
}

func (p ServerPacketSetCompression) Write(b common.Buffer) error {
	return b.WriteVarInt(p.Threshold)
}

func (p ServerPacketSetCompression) String() string {
	return fmt.Sprintf("ServerPacketSetCompression{Threshold: %d}", p.Threshold)
}
