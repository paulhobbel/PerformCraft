package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
)

type ClientPacketTeleportConfirm struct {
	TeleportID int32
}

func (ClientPacketTeleportConfirm) ID() base.PacketID {
	return PlayTeleportConfirm
}

func (p *ClientPacketTeleportConfirm) Read(buffer bufio.ByteBuffer) (err error) {
	p.TeleportID, err = buffer.ReadVarInt()

	return
}

func (p ClientPacketTeleportConfirm) Write(buffer bufio.ByteBuffer) error {
	return buffer.WriteVarInt(p.TeleportID)
}

func (p ClientPacketTeleportConfirm) String() string {
	return fmt.Sprintf("ClientPacketTeleportConfirm{TeleportID: %d}", p.TeleportID)
}

type ChatMode int32

const (
	ChatEnabled ChatMode = iota
	ChatDisabled
)

func (m ChatMode) String() string {
	if m == ChatEnabled {
		return "Enabled"
	}

	return "Disabled"
}

type ClientPacketSettings struct {
	Locale       string
	ViewDistance byte
	ChatMode     ChatMode
	ChatColors   bool
	SkinParts    byte
	MainHand     int32
}

func (ClientPacketSettings) ID() base.PacketID {
	return PlayClientSettings
}

func (p *ClientPacketSettings) Read(buffer bufio.ByteBuffer) (err error) {
	p.Locale, err = buffer.ReadString()
	p.ViewDistance, err = buffer.ReadByte()

	chatMode, err := buffer.ReadVarInt()
	p.ChatMode = ChatMode(chatMode)

	p.ChatColors, err = buffer.ReadBoolean()
	p.SkinParts, err = buffer.ReadByte()
	p.MainHand, err = buffer.ReadVarInt()

	return
}

func (p ClientPacketSettings) Write(buffer bufio.ByteBuffer) (err error) {
	err = buffer.WriteString(p.Locale)
	err = buffer.WriteByte(p.ViewDistance)
	err = buffer.WriteVarInt(int32(p.ChatMode))
	err = buffer.WriteBoolean(p.ChatColors)
	err = buffer.WriteByte(p.SkinParts)
	err = buffer.WriteVarInt(p.MainHand)

	return
}

func (p ClientPacketSettings) String() string {
	return fmt.Sprintf("ClientPacketSettings{Locale: %s, ViewDistance: %d, ChatMode: %s, ChatColors: %v, SkinParts: %d, MainHand: %d}",
		p.Locale, p.ViewDistance, p.ChatMode, p.ChatColors, p.SkinParts, p.MainHand)
}

type ClientPacketPluginMessage struct {
	Channel string
	Data    []byte
}

func (ClientPacketPluginMessage) ID() base.PacketID {
	return PlayPluginMessageClient
}

func (p *ClientPacketPluginMessage) Read(buffer bufio.ByteBuffer) (err error) {
	p.Channel, err = buffer.ReadString()
	p.Data = buffer.Bytes()

	return
}

func (p ClientPacketPluginMessage) Write(buffer bufio.ByteBuffer) (err error) {
	err = buffer.WriteString(p.Channel)
	_, err = buffer.Write(p.Data)

	return
}

func (p ClientPacketPluginMessage) String() string {
	return fmt.Sprintf("ClientPacketPluginMessage{Channel: %s, Data: %v}", p.Channel, p.Data)
}

type ClientPacketKeepAlive struct {
	KeepAliveID int64
}

func (ClientPacketKeepAlive) ID() base.PacketID {
	return PlayKeepAliveClient
}

func (p *ClientPacketKeepAlive) Read(buffer bufio.ByteBuffer) (err error) {
	p.KeepAliveID, err = buffer.ReadLong()

	return
}

func (p ClientPacketKeepAlive) Write(buffer bufio.ByteBuffer) error {
	return buffer.WriteLong(p.KeepAliveID)
}
